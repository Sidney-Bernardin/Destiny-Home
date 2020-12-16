# Destiny Home


### What is this?
Destiny Home is the backend to my Google Assistant action. 

##### My Google Home Action
My Google Home action is a Destiny Item Manager (DIM) that I created with [Actions onGoogle](https://console.actions.google.com). A Destiny Item Manager can do things like move and equip item within the game of Destiny. [Here is an example of a popular third party Destiny Item Manager](https://app.destinyitemmanager.com/).  

##### My Webhook / Google Cloud Function
For example after my Google Home action recives a request from a client to equip an armor-pice, it will make an HTTP request to werever the Webhook / Google Cloud Function is living on the WWW, and tell it to equip that item.
