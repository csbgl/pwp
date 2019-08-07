# pwp

You install a software that requires a password? Simple you provide that.

But what happens if you want to daemonize or automate that? 

- Create a Password file in plaintext? No
- Create a Password file and encrypt that? Better.
  - Use symmentric encryption? No. You need to store the PSK in plain.
  - Use asymetric encryption? No. You need to have the privekey in plaintext.

That's what PWP can help.
