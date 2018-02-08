# hass-go

Home-Assistant external go process for weather, calendar, season, sensors, light.

This process will read data from:
- DarkSky for the weather, current and hourly:
  - Clouds
  - Rain
  - Temperature
- iCloud for calendars
  - School, Free, Work
  - Season: Winter, Summer, Autumn, Spring
  - Time-Of-Day: breakfast, morning, noon, afternoon, evening, bedtime, sleeptime, night
  - Weather: When to send a weather report with a location (e.g. Asia/Shangai) and time-window.
  - Lights: Brightness and Temperature modifiers

The Calendars will create the following HTTP sensors:
- sensor.jennifer (school, free)
- sensor.sophia (school, free)
- sensor.parents (work, free)

Others are:
- sensor.time_of_day

A computation is done for the supplied location (lat, long) using weather, sun position and season to create values for the following HTTP sensors:
- sensor.lights_hue_ct
- sensor.lights_hue_bri

And it will do the following:
- Post a json message to a running home-assistant instance http://IP:Port/api/states/sensor.NAME to update HTTP sensors
- Post messages to slack for
  - Weather report
  - AQI (Air Quality Index)
  - Notifications (configurable), these are examples:
    - rain
    - freezing

