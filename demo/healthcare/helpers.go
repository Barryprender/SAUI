package healthcare

// Doctor is a practitioner entry.
type Doctor struct {
	Slug        string
	Name        string
	Specialty   string
	AvatarClass string
}

// SlotDef is a fixed appointment slot in the schedule.
type SlotDef struct {
	ID       string
	Doctor   Doctor
	Day      string
	DayOrder int
	Time     string
}

var allDoctors = []Doctor{
	{Slug: "patel", Name: "Dr. Priya Patel", Specialty: "General Practice", AvatarClass: "avatar--violet"},
	{Slug: "clarke", Name: "Dr. James Clarke", Specialty: "Cardiology", AvatarClass: "avatar--blue"},
}

var allDays = []struct {
	Name  string
	Slug  string
	Order int
}{
	{"Monday", "mon", 1},
	{"Wednesday", "wed", 2},
	{"Friday", "fri", 3},
}

var allTimes = []struct {
	Display string
	Slug    string
}{
	{"09:00", "0900"},
	{"14:30", "1430"},
}

var allSlots []SlotDef

func init() {
	for _, day := range allDays {
		for _, t := range allTimes {
			for _, doc := range allDoctors {
				allSlots = append(allSlots, SlotDef{
					ID:       doc.Slug + "-" + day.Slug + "-" + t.Slug,
					Doctor:   doc,
					Day:      day.Name,
					DayOrder: day.Order,
					Time:     t.Display,
				})
			}
		}
	}
}

func AllSlots() []SlotDef { return allSlots }

func SlotByID(id string) (SlotDef, bool) {
	for _, s := range allSlots {
		if s.ID == id {
			return s, true
		}
	}
	return SlotDef{}, false
}
