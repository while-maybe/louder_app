package domain

type GenTag string // General Tags: Legendary, Top 10,
type GenreTag string
type MoodTag string

type Album struct {
	Title       string
	Artist      Entity   // A Band or a leading artist: U2, Keith Jarrett
	SuppArtists []Entity // Supporting Artists if any
	Setting     string   // Studio, Live, Compilation
	Location    string   // Wembley stadium, etc
	Description string
	GenTags     []GenTag
	MainGenre   string     // Classical music, Jazz, Rock, Metal
	SubGenres   []GenreTag // Pop-Rock, Be-Bop, Melodic Metal
	MoodTags    []MoodTag  // Intimate, Exciting, Nostalgic
	Composer    Entity     // if any
}

type Release struct {
	Album    Album
	Version  string // Original, Remastered, Extended, etc
	Medium   string // CD, Vinyl, K7, Digital, MD, 8-track...
	Year     byte
	Discs    []Disc
	Duration uint
	CoverURL string
	CoR      string // Country of Release - Japan has amazing releases
	Label    Label

	// Below for classical music ?
	Conductor Entity // if any
	Band      Entity // if any
}

type Disc struct {
	Duration uint
	Tracks   []Track
}

type Track struct {
	Number   byte
	Name     string
	Duration uint
	Lyrics   string
}

type Label struct {
	Name       Entity  // Sony Music, ECM Records
	Previously []Label // In case they were acquired by other labels
}

type Entity struct {
	Name        string // Name of Entity/Band
	Description string // Description of the Entity, Band, Compose, etc
	CoO         string // Country of Origin
	YStart      byte   // Year the Entity/Band started
	YEnd        byte   // Year the Entity/Band ended
	ImageURLs   []string
	EntityURL   string
}
