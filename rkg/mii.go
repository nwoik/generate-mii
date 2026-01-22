package rkg

type Mii struct {
	Name     string
	Creator  string
	Gender   bool // false = male, true = female
	Favorite bool
	Height   int
	Weight   int

	FaceShape int
	SkinColor int

	HairType  int
	HairColor int
	HairFlip  bool

	EyeType     int
	EyeColor    int
	EyeSize     int
	EyeY        int
	EyeX        int
	EyeRotation int

	NoseType int
	NoseSize int
	NoseY    int

	MouthType  int
	MouthColor int
	MouthSize  int
	MouthY     int

	HasGlasses   bool
	GlassesType  int
	GlassesColor int
	GlassesY     int
}
