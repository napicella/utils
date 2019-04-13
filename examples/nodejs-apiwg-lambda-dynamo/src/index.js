/**
 *
 * Simple lambda which converts markddown to html and stores the result in DyanmoDB.
 * It provides two inner endpoints /get and /store.
 * The request is routed to the correct handler with aws-lambda-router.
 *
 */

/**
 * Provides a router for lambdas which uses APi gateway proxy integration and {proxy +}
 * In this mode, Api gateway just forwards the http request to the lambda which makes it available to the handler in the
 * event object.
 * aws-lambda-router routes the request based on the url in the event object.
 * @see https://github.com/spring-media/aws-lambda-router
 */
const router = require('aws-lambda-router');
/**
 * markdown to html
 * @see https://github.com/showdownjs/showdown
 */
const showdown  = require('showdown'),
  converter = new showdown.Converter({simpleLineBreaks : false, noHeaderId : true});

const AWS = require('aws-sdk');
const changeLogTableName = process.env.CHANGE_LOG_TABLE;
const dynamodb = new AWS.DynamoDB({apiVersion: '2012-08-10'});


// handler for an api gateway event
exports.handler = router.handler({
  proxyIntegration: {
    debug: true,
    routes: [
      {
        path: '/get/:title',
        method: 'GET',
        action: (request, context) => {
          return get(request.paths.title);
        }
      },
      {
        path: '/store/:title',
        method: 'POST',
        // we can use the path param 'title' in the action call:
        action: (request, context) => {
          let title = request.paths.title;
          let md = request.body.markdown;
          let html = converter.makeHtml(md);
          return store(title, html, md);
        }
      }
    ]
  }
});

function store(title, html, md) {
  return dynamodb.putItem({
    "TableName": changeLogTableName,
    "Item": {
      "title"  : {"S": title},
      "date"   : {"S" : new Date().toISOString() },
      "html"   : {"S": html},
      "md"     : {"S": md}
    }
  }).promise();
}

function get(title) {
  let params = {
    ExpressionAttributeNames: {
      "#t": "title"
    },
    ExpressionAttributeValues: {
      ":v1": {
        S: title
      }
    },
    KeyConditionExpression: "#t = :v1",
    TableName: changeLogTableName
  };
  return dynamodb.query(params).promise()
                 .then(result => {
                   if (result.Items && result.Items.length > 0) {
                     return { result : result.Items[0] };
                   }

                   return {};
                 }).catch(e => e);
}
