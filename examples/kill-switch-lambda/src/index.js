const AWS = require('aws-sdk'); // eslint-disable-line import/no-extraneous-dependencies

// https://docs.aws.amazon.com/cli/latest/reference/apigateway/update-stage.html
const apigateway = new AWS.APIGateway({apiVersion: '2015-07-09', region: process.env.AWS_REGION});
var params = {
    restApiId: process.env.API_GW_ID,
    stageName: process.env.API_GW_STAGE,
    patchOperations: [{
        op: 'replace',
        path: '/~1*/*/throttling/rateLimit',
        value: '0'
    }]
};
console.log(`Kill switch for ${params.restApiId} - ${params.stageName}`);

const handler = (event, context, callback) => {

    apigateway.updateStage(params, function(err, data) {
        if (err) {
            console.log(err, err.stack);
            callback(err);
        } else {
            console.log("All the request will be throttled");
            callback(null, {
                statusCode: 200,
                body: JSON.stringify({message : "All the request will be throttled!!"})
            });
        }
    });

}

module.exports = {
    handler,
};
