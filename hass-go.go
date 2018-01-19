package main

// This can be computed from sun-rise, sun-set
// - sensor.dark_or_light (Dark, Twilight, Light)

// The data for these sensors can also come from a calendar
// In icloud we have the following calendars already:
// - TimeOfDay
// - Season
// - Weather to report
// - Lights (NA)
// Note: These calendars can be downloaded on a low frequency.
// - Weather (The moments when to report the weather)
// - TimeOfDay (Breakfast, Morning, Lunch, Afternoon, Evening, ..)
// - Season (Summer, Winter, Spring, Autumn)
// - Lights: Temperature (min(154) - max(500)) and Brightness (0.0 - 1.0)

// Every X minutes:
// - Download calendars from icloud URL
// - Parse all calendars
// - For 'Today' and 'Tomorrow' get all events
// - Search events for:
//   - "Jennifer.School" & "Jennifer.Free"
//   - "Sophia.School" & "Sophia.Free"
//   - "Parents.Work" & "Parents.Free"
//   - "Alarm.Event"
// - From the events build the info for all the sensors
// - Today school time is from 6:00 AM to 17:00 AM
//   So if we have school tomorrow then the sensor will still be set
//   to School otherwise it would be Free
// - Push the sensor state of Jennifer, Sophia, Parents and Alarm to HASS
// - Sleep

func main() {

}
