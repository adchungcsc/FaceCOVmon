import * as cdk from "@aws-cdk/core";
import * as dynamodb from "@aws-cdk/aws-dynamodb"
import * as lambda from "@aws-cdk/aws-lambda"
import * as apigateway from "@aws-cdk/aws-apigateway";

export class ApiStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const observationsDatabase = new dynamodb.Table(this, "Observations",{
            partitionKey: { name: "cameraID", type: dynamodb.AttributeType.STRING },
            sortKey: { name: "date", type: dynamodb.AttributeType.STRING },
            tableName: "Observations",
        });

        //Lambda func. to handle the incoming observation
        const handler = new lambda.Function(this, "FaceCOVHandler", {
            runtime: lambda.Runtime.GO_1_X,
            memorySize: 256,
            functionName: "FaceCOVHandlerAPI",
            code: lambda.Code.fromAsset("../backend/api/function/handler.zip"),
            handler: "handler",
        });

        //Let the lambda API have RW access to the db
        observationsDatabase.grantReadWriteData(handler);

        const api = new apigateway.RestApi(this, "FaceCOVMonAPI");
        const upload = api.root.resourceForPath("/data");
        upload.addMethod("GET", new apigateway.LambdaIntegration(handler));
        const data = api.root.resourceForPath("/upload");
        data.addMethod("PUT", new apigateway.LambdaIntegration(handler));
    }
}

