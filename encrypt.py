import pyAesCrypt
# encryption/decryption buffer size - 64K
bufferSize = 64 * 1024
password = input("Enter password: ")
# decrypt
pyAesCrypt.encryptFile("hass-go-calendar.TOML", "hass-go-calendar.TOML.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-sensors.TOML",  "hass-go-sensors.TOML.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-slack.TOML",    "hass-go-slack.TOML.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-weather.TOML",  "hass-go-weather.TOML.aes", password, bufferSize)
pyAesCrypt.encryptFile("hass-go-secrets.TOML",  "hass-go-secrets.TOML.aes", password, bufferSize)
