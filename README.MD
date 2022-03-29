# Boox bot
A telegram bot for sending books to your Boox device.You may try it on [Telegram](https://t.me/Boox_sync_bot)

Currently,it only supports to login [China](https://send2boox.com/) server.

## Feature

* Login to boox
* send book to boox
* support in different server,e.g. China/US/Vietnam

## How it works

I find out that the boox server using JWT token to authenticate after login,if I can get the token then I can use it to upload your books to the boox server and push it to your server.
This is the starting point of this bot.

### Login

The /code command uses the `users/sendMobileCode` endpoint to you an email with a login code which the bot will use to login to boot,obtain a token by the endpoint `users/signupByPhoneOrEmail`.

The token will expire in six months,as far I know.

### Upload

Once the bot receive a document,it will first get a temporary aliyun [sts](https://help.aliyun.com/document_detail/410749.html) by the endpoint `config/stss`, 
then use it to upload to Aliyun OSS under the path of your uid.Everything you upload is only accessibly by you.


## Can this bot access my other information

In theory,it has obtained your access token which means it can do everything just like you are in send2boox.com.
But this bot only keep your token to upload documents and push it to your device,nothing more.