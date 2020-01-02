### Poor's man kill switch for Api Gateway
Given a topic (like the one created by a CLoudWatch alarm) and an Api gateway Id
registers a Lambda that is going to change the Api gateway global throttling
configuration to throttle all the requests.

#### Why?
It allows to expose a public endpoint for a live demo, while giving you the peace
of mind that in case of abuse, the api will kill the requests keeping the cost very low.
This is possible thanks to API gateway cost model: throttle requested are not
considered in the billing :)

