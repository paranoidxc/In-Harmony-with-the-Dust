/*
anything time you use a machine learning model
you want to ask youself the question

what data was used to train this model
who trained this model
what context and for what reason who trained this model
*/

let features;
let video;
let knn;
let labelP;
let ready = false;

function modelReady() {
  console.log("Model is ready");
}

function videoReady() {
  console.log("Video is ready");
}

function setup() {
  background(0);
  createCanvas(640, 560);
  video = createCapture(VIDEO);
  video.hide();
  features = ml5.featureExtractor("MobileNet", modelReady);
  knn = ml5.KNNClassifier();
  labelP = createP("need training data");
}

function goClassify() {
  const logits = features.infer(video);
  knn.classify(logits, function (error, result) {
    if (error) {
      console.error(error);
    } else {
      //console.log(result);
      labelP.html(result.label);
      goClassify();
    }
  });
}

/*
function mousePressed() {
  if (knn.getNumLabels() > 0) {
    const logits = features.infer(video);
    knn.classify(logits, gotResult);
  }
}
*/

function keyPressed() {
  const logits = features.infer(video);
  if (key == "l") {
    knn.addExample(logits, "left");
    console.log("left");
  } else if (key == "r") {
    knn.addExample(logits, "right");
    console.log("right");
  } else if (key == "u") {
    knn.addExample(logits, "up");
    console.log("up");
  } else if (key == "d") {
    knn.addExample(logits, "down");
    console.log("down");
  } else if (key == "s") {
    knn.save("model.json");
  }
  //console.log(logits.dataSync());
}

function draw() {
  background(0);
  image(video, 0, 0);

  if (!ready && knn.getNumLabels() > 0) {
    goClassify();
    ready = true;
  }
}
