import pyAesCrypt
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
password = password = input("Enter password: ")
# decrypt
pyAesCrypt.decryptFile("hass-go-calendar.toml.aes", "hass-go-calendar.toml", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-sensors.toml.aes", "hass-go-sensors.toml", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-slack.toml.aes", "hass-go-slack.toml", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-weather.toml.aes", "hass-go-weather.toml", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-secrets.toml.aes", "hass-go-secrets.toml", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-suncalc.toml.aes",  "hass-go-suncalc.toml", password, bufferSize)
pyAesCrypt.decryptFile("hass-go-lighting.toml.aes",  "hass-go-lighting.toml", password, bufferSize)
