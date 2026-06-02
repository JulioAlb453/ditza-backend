package domain

import "fmt"

type Mood string

const (
	MoodHappy    Mood = "happy"
	MoodNeutral  Mood = "neutral"
	MoodSad      Mood = "sad"
	MoodSleeping Mood = "sleeping"
)

func ParseMood(raw string) (Mood, error) {
	mood := Mood(raw)
	switch mood {
	case MoodHappy, MoodNeutral, MoodSad, MoodSleeping:
		return mood, nil
	default:
		return "", fmt.Errorf("estado de ánimo de mascota inválido: %s", raw)
	}
}

func (m Mood) String() string {
	return string(m)
}

func (m Mood) IsValid() bool {
	_, err := ParseMood(string(m))
	return err == nil
}
