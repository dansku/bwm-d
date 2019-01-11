# Bandwidth Monitoring Daemon

## Requirements

  * A SBC Supported by Golang and [balena.io](https://www.balena.io) (and obviously an account)
  * An (industrial) account on [ubidots](https://ubidots.com/)
  
## Usage

  * Create an account on both balena and ubidots.
  * Create a bwm-d app on balena
    * Under `Environment Varialbes`, set your Ubidots Token
  * Clone it locally and push into your balena repo
  * Create Devices and download their boot images. Burn, boot your pi and start playing around :)

## Shortcomings

  * I can't run `dep` on balena somehow, thus why I'm adding the `vendor` directory