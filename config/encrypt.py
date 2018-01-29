import pyAesCrypt
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
password = input("Enter password: ")
# decrypt
pyAesCrypt.encryptFile("hass-go-calendar.toml", "hass-go-calendar.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-sensors.toml",  "hass-go-sensors.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-slack.toml",    "hass-go-slack.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-weather.toml",  "hass-go-weather.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-secrets.toml",  "hass-go-secrets.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-suncalc.toml",  "hass-go-suncalc.toml.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-lighting.toml",  "hass-go-lighting.toml.aes", password, bufferSize)
