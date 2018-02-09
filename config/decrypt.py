import pyAesCrypt
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
password = password = input("Enter password: ")
# decrypt
pyAesCrypt.decryptFile("aqi.json.aes", "aqi.json", password, bufferSize)
pyAesCrypt.decryptFile("calendar.json.aes", "calendar.json", password, bufferSize)
pyAesCrypt.decryptFile("hass.toml.aes", "hass.toml", password, bufferSize)
pyAesCrypt.decryptFile("lighting.toml.aes", "lighting.toml", password, bufferSize)
pyAesCrypt.decryptFile("reporter.json.aes", "reporter.json", password, bufferSize)
pyAesCrypt.decryptFile("secrets.toml.aes", "secrets.toml", password, bufferSize)
pyAesCrypt.decryptFile("sensors.toml.aes", "sensors.toml", password, bufferSize)
pyAesCrypt.decryptFile("shout.toml.aes", "shout.toml", password, bufferSize)
pyAesCrypt.decryptFile("suncalc.toml.aes", "suncalc.toml", password, bufferSize)
pyAesCrypt.decryptFile("timeofday.json.aes", "timeofday.toml", password, bufferSize)
pyAesCrypt.decryptFile("weather.toml.aes", "weather.toml", password, bufferSize)
