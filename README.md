# Bandwidth Monitoring Daemon

## Prerequisites

Requires:

  * balena.io
  * ubidots
  
## Usage

  * Create an account on both balena and ubidots.
  * Create a bwm-d app on balena
    * Under Application | Environment Varialbes, set your Ubidots Token
  * Clone it locally and push into your balena repo
  * Download Raspberry Pi Images for your app and start playing around :)

## Shortcomings

  * I can't run `dep` on balena somehow, thus why I'm adding the `vendor` directory