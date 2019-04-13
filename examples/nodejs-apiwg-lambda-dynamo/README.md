Api gateway proxy integration to a Node Js Lambda.  
The lambda converts markdown and store the result in dyanmo db.  
It provides two inner endpoints /get and /store.  
The request is routed to the correct handler with aws-lambda-router.  
