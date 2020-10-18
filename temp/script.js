// 1). grab needed elements' reference
// #video element
const video = document.getElementById('video');

// 2). load the models
async function loadModels() {
    /* ssdMobilenetv1 - for face detection */
    await faceapi.nets.ssdMobilenetv1.loadFromUri('./models');

    // start the live webcam video stream
    startVideoStream();
}

// 3). startVideoStream() function definition
function startVideoStream() {
    // check if getUserMedia API is supported
    if (navigator.mediaDevices.getUserMedia) {
        // get the webcam's video stream
        navigator.mediaDevices.getUserMedia({ video: true })
            .then(function(stream) {
                // set the video.srcObject to stream
                video.srcObject = stream;
            })
            .then(makePredictions)
            .catch(function(error) {
                /* if something went wrong, just console.log() error. */
                console.log(error);
            });
    }
}


// 4). makePredictions() function definition
function makePredictions() {
    // get the #canvas
    const canvas = document.getElementById('canvas');
    // resize the canvas to the #video dimensions
    const displaySize = { width: video.width, height: video.height };
    faceapi.matchDimensions(canvas, displaySize);

    /* get "detections" for every 500 milliseconds */
    setInterval(async function() {
        /* this "detections" array has all the things like the "prediction results" as well as the "bounding box" configurations! */
        const detections = await faceapi.detectAllFaces(video);
        document.getElementById("detected").innerHTML = ("found " + detections.length + " individuals");
        //If face detected, post to the backend for processing.
        if(detections.length > 0){

            //conver to a png to send to the backend API
            var dataURL = canvas.toDataURL('image/png);
            fetch('https://api.mocki.io/v1/ddb6d39b', { method: 'POST', body: {} })
                .then(response => response.json())
                .then(data => console.log(data);
        }
        /* resize the detected boxes to match our video dimensions */
        const resizedDetections = faceapi.resizeResults(detections, displaySize);


        // before start drawing, clear the canvas
        canvas.getContext('2d').clearRect(0, 0, canvas.width, canvas.height);
        // use faceapi.draw to draw "detections"
        faceapi.draw.drawDetections(canvas, resizedDetections);

        // to draw expressions
        //faceapi.draw.drawFaceExpressions(canvas, resizedDetections);
        // to draw face landmarks
        //faceapi.draw.drawFaceLandmarks(canvas, resizedDetections);
    }, 500);
}




// activate the loadModels() function
loadModels()

