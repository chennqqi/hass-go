import pyAesCrypt
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
password = input("Enter password: ")
# decrypt
pyAesCrypt.encryptFile("aqi.json",      "aqi.json.aes", password, bufferSize)
pyAesCrypt.encryptFile("calendar.json", "calendar.json.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass.toml",     "hass.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("lighting.toml", "lighting.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("reporter.json", "reporter.json.aes", password, bufferSize)
pyAesCrypt.encryptFile("secrets.toml",  "secrets.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("sensors.toml",  "sensors.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("shout.toml",    "shout.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("suncalc.toml",  "suncalc.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("timeofday.json","timeofday.json.aes", password, bufferSize)
pyAesCrypt.encryptFile("weather.toml",  "weather.toml.aes", password, bufferSize)
