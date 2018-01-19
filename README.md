# hass-go

Home-Assistant external go process for weather, calendar, season, sensors, light.

This process will read data from:
- DarkSky for the weather
- iCloud for calendars
  - School, Free, Work
  - Season: Winter, Summer, Autumn, Spring
  - Time-Of-Day: breakfast, morning, noon, afternoon, evening, bedtime, sleeptime, night
  - Weather: When to send a weather report with a location (e.g. Asia/Shangai) and time-window.
  - Lights: Brightness and Temperature modifiers

A computation is done for the supplied location (lat, long) to create a sensor that can have the following cycling states:
- Night
- Astronomical twilight
- Nautical twilight
- Civil twilight
- Sunrise
- Morning
- Noon
- Afternoon
- Sunset
- Civil twilight
- Nautical twilight
- Astronomical twilight
- Night

These states are send to home-assistant as a 'sensor.


And it will do the following:
- Post a json message to a running home-assistant instance http://IP:Port/api/states/sensor.NAME to update HTTP sensors
- Post messages to slack for
  - Weather report
  - Notifications (configurable), these are examples:
    - rain
    - freezing

