That is an excellent and critical question. It lies at the heart of the trade-offs between idiomatic Go and classical object-oriented encapsulation.

You are correct: **exporting fields *does* give the consumer the technical ability to change them.**

However, in this specific architectural pattern (an immutable data record created by a trusted source), exporting the fields is the **standard, idiomatic, and recommended practice in Go.**

Let's break down exactly why this is considered safe and preferable in this context.

### The "Social Contract" vs. The "Technical Contract"

*   **Technical Contract:** What the compiler enforces. Can a consumer write `myDiceRoll.RollSum = 999`? With an exported field, **yes, they can.**
*   **Social Contract:** How developers agree to use an API. In Go, the "social contract" for a struct with no setters is that it's intended to be **read-only**.

Go's philosophy leans heavily on the social contract, preferring clarity and simplicity over the rigid technical enforcement found in languages like Java or C#. The idea is that if an object is clearly designed to be an immutable data carrier, other developers on the team (or users of your library) will respect that design.

### Why It's Safe in *Your* Architecture

The safety of this pattern depends entirely on the "gatekeeper" you've already built. Think about the flow of data:

1.  **The Untrusted World:** An HTTP request comes in with user input (e.g., `numDice=50`).
2.  **The Gatekeeper (Adapter):** Your `StdlibGenerator.GenerateDiceRoll` method receives this input.
    *   It **validates** the input (`50` is invalid).
    *   It **generates** a valid result using `math/rand`.
    *   It **creates** a `domain.RandomDice` struct, which is now a snapshot of a valid, completed event.
3.  **The Trusted World:** Inside your `Service` and `Handler`, you are now holding this `RandomDice` object. **At this point, its job is done.** It is simply a container for results.

The critical insight is: **Why would any code inside your service *need* to change the `RollSum` of a roll that has already happened?**

There is no legitimate business reason to do so. The roll is a historical fact. Changing `RollSum` would be a bug. The Go philosophy trusts the developer not to introduce such a bug deliberately.

### The Trade-Off: Simplicity vs. Absolute Prevention

Let's compare the two approaches directly.

#### Approach 1: Exported Fields (Idiomatic Go)

```go
// in domain/random.go
type RandomDice struct {
	Roll    []uint
	RollSum uint
}

// in service/game_service.go
// Code is clean and direct.
total := myDiceRoll.RollSum
```
*   **Pros:**
    *   **Simple & Clear:** `myDiceRoll.RollSum` is more readable than `myDiceRoll.GetSum()`.
    *   **Less Boilerplate:** No need to write getter methods for every field.
*   **Cons:**
    *   **Technically Mutable:** Another developer could maliciously or accidentally write `myDiceRoll.RollSum = 999`.

#### Approach 2: Unexported Fields + Getters (Classic OOP)

```go
// in domain/random.go
type RandomDice struct {
	roll    []uint
	rollSum uint
}

func (rd *RandomDice) Sum() uint { return rd.rollSum }
// We also need a getter for the slice, with a defensive copy!
func (rd *RandomDice) Roll() []uint {
    //... copy logic ...
}

// in service/game_service.go
// Code is slightly more verbose.
total := myDiceRoll.Sum()
```
*   **Pros:**
    *   **Technically Immutable:** The compiler will prevent `myDiceRoll.rollSum = 999`. The object is truly sealed.
*   **Cons:**
    *   **Verbose:** Requires writing (and maintaining) a getter for every single field.
    *   **Considered "Stutter":** Go developers often see this as unnecessary "syntactic stutter" when the getter adds no logic.

### When Unexported Fields are the Right Choice

The Go community absolutely uses unexported fields and getters, but reserves them for situations where they provide real value:

1.  **When the struct is MUTABLE:** If you have `Set...` methods, the fields **must** be unexported to force all changes to go through the validating setter.
2.  **When a "getter" computes a value:** `rect.Area()` is a perfect example. The area isn't a stored field.
3.  **When enforcing concurrency safety:** A getter can wrap a mutex lock around reading a value.

Your `RandomDice` struct fits none of these criteria. It's an immutable data transfer object.

### Conclusion: Trust the Pattern

You have designed a robust system where **creation is strictly controlled**. Because of this, you can afford to make **reading the data simple and direct**.

**The final, idiomatic recommendation remains:** **Use exported fields for this type of immutable data record.**

It feels counter-intuitive if you come from a background where private fields and public getters are the default for everything. But in the context of Go's design philosophy and the specific architecture you've built, it is the standard and preferred way. You are trading a small amount of theoretical risk (that a developer will do something they shouldn't) for a large gain in code clarity and simplicity.