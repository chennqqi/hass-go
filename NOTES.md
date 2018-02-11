# task graph with dependency and property changed triggers

## aqi has the following triggers

- property aqi.tick

## calendar has the following triggers

- property calendar.tick

## hass has the following triggers

NOTE: Maybe a hass.sensor makes more sense

- anything written to a 'sensor.NAME' property

## lighting has the following triggers

- property weather.currently:clouds
- property time.season
- property lighting.tick

## reporter has the following triggers

- anything written to the 'reporter' domain

## sensors has the following triggers

NOTE: Maybe a sensors.sensor makes more sense to be triggered

- all properties as listed in the configuration

## suncalc has the following triggers

- property time.day

## timeofday has the following triggers

- property timeofday.tick

## weather has the following triggers

- property weather.tick
