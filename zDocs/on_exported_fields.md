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

You have hit upon one of the most important and defining questions in Hexagonal Architecture. This is where the discipline really pays off.

The strong and widely-accepted best practice is to **add a mapper in the adapter.**

**Do NOT put `json` tags (or `db`, `xml`, etc.) in your core domain entity.**

Let's explore exactly why this is so critical.

---

### The Core Principle: Separation of Concerns

The entire point of the hexagon is to protect your core domain logic from the outside world.

*   **Your Domain (`internal/domain`)** should know *only* about business rules. It should have no knowledge of how it is stored, how it is displayed, or how it is transmitted over a network. Is it JSON? XML? A command-line printout? The domain doesn't know and shouldn't care.
*   **Your Adapter (`internal/adapters/http`)** is a "translator." Its entire job is to translate between the "language" of the outside world (e.g., HTTP requests and JSON) and the "language" of your domain (your pure domain objects and methods).

Putting a `json:"roll_sum"` tag on your `RandomDice` domain entity creates a **leaky abstraction**. Your pure business logic now has a dependency on a specific presentation detail (`json`). This violates the fundamental **Dependency Rule**: dependencies must always point inwards, from adapters to the core. A `json` tag creates a conceptual dependency pointing outwards from the core to the adapter's technology.

### The Correct Approach: DTOs and Mappers in the Adapter

Here is the clean, hexagonal way to solve this problem.

1.  **Keep your domain object pure.** It has no tags.
2.  **Create a DTO (Data Transfer Object)** in your `http` adapter package. This struct's only purpose is to define the shape of your JSON response. *This* is where the JSON tags belong.
3.  **Create a mapper function** that converts your domain object into your DTO.

Let's see it in code.

**1. Your Pure Domain Object (No Changes)**

```go
// internal/domain/dice.go
package domain

// ... (NewRandomDice, etc.)

type RandomDice struct {
	roll    []uint
	rollSum uint
}

func (d *RandomDice) Roll() []uint { /* ... returns copy ... */ }
func (d *RandomDice) Sum() uint { return d.rollSum }
```
Notice: No `json` tags. Perfect.

**2. The DTO and Mapper in the Adapter**

```go
// internal/adapters/http/dto.go
package http

import "my-app/internal/domain"

// DiceResponse is a DTO specifically for the JSON API response.
// It's a "bag of data" with presentation tags.
type DiceResponse struct {
	Rolls []uint `json:"rolls"` // Control the JSON field name
	Sum   uint   `json:"sum"`
}

// toDiceResponse is a mapper function. It maps a domain object to a DTO.
// It lives in the adapter because it's part of the adapter's translation logic.
func toDiceResponse(d *domain.RandomDice) DiceResponse {
	return DiceResponse{
		Rolls: d.Roll(), // Uses the public getter method
		Sum:   d.Sum(),    // Uses the public getter method
	}
}
```

**3. How it's Used in the Handler**

```go
// internal/adapters/http/handler.go
package http

func (h *DiceHandler) HandleRollDice(w http.ResponseWriter, r *http.Request) {
	// 1. Get parameters from request...
	// 2. Call your application service, which returns the rich domain object.
	diceRoll, err := h.diceService.Roll(2, 6) // Returns a *domain.RandomDice
	if err != nil {
		// ... handle error
		return
	}

	// 3. Use the mapper to convert the domain object to the response DTO.
	response := toDiceResponse(diceRoll)

	// 4. Encode the DTO. The json tags in the DTO will be used.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response) // Produces {"rolls":[...],"sum":...}
}
```

### Why is this better? The Long-Term Benefits

| Aspect | Mapper in Adapter (Best Practice) | Tags in Domain (Anti-Pattern) |
| :--- | :--- | :--- |
| **Separation** | **Excellent.** Domain knows nothing about JSON. Adapter handles all translation. | **Poor.** Domain is now coupled to JSON. It's a "leaky abstraction." |
| **Flexibility** | **High.** Want to add a gRPC endpoint? Create a new `grpc` adapter with its own protobuf-based DTOs. The domain doesn't change. | **Low.** If you add a gRPC adapter, do you add `protobuf` tags to the domain object? What about `xml` tags? The domain becomes a dumping ground for presentation details. |
| **Maintainability** | **Clean.** Each adapter is self-contained. It's clear where translation logic lives. | **Messy.** Changes to the API format might require changing the core domain object, which is risky and counter-intuitive. |
| **Boilerplate** | **Slightly more.** You have to write the DTO struct and the mapper function. | **Less.** It's quicker to just add a tag. |

The slightly increased boilerplate of the mapper approach is the price you pay for a clean, maintainable, and truly decoupled architecture. It is a trade-off that is **almost always worth it** in any project that is expected to grow or be maintained over time.

