# Using REDIS and splitting all modules into micro-services

Why ?

- Every micro-service can wait for an external signal to update itself or
  it can run at a certain time-period.
  This could be done with Pub/Sub and having a dedicated micro-service
  running that sends 'tick' messages at certain periods to micro-services 
  that are running on a certain frequency.
- Easier to debug and if the micro-service crashes it can restart
- State can be saved into REDIS and every micro-service can read other 
  micro-service properties from Redis.
- Configuration of every micro-service can be pushed into REDIS and every
  micro-service could restart itself when it detects that the configuration 
  has changed.


  The modules are:
  
  - AQI
    - Update = Tick
    - HTTP GET from designated URL and parse JSON
    - Write 'weather.currently:aqi' = 0-1000 value to REDIS
  
  - CALENDAR
    - Update = Tick
    - Read calendar information from REDIS
    - Write 'weather.currently:aqi' = 0-1000 value to REDIS
  
  - CALENDAR CHANGED
    - Update = Tick 5 minutes
    - Will HTTP GET calendars and post them into REDIS
    - Will send tick to CALENDAR when any calendar has changed (content hash)

  - HASS
    - Update = Tick 1 minute
    - Get information from REDIS and HTTP POST to Home-Assistant

  - LIGHTING
    - Update = Tick 1 minute
  
  - REPORTER
  
  - SENSORS
  
  - SHOUT
  
  - SUNCALC
    - Update = Tick 1 hour
  
  - TIMEOFDAY
    - Update = Tick every 15 minutes of the hour
  
  - WEATHER
    - Update = Tick every 15 minutes of the hour
    - Push information to REDIS

