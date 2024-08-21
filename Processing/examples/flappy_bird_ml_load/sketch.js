let birds = [];
let pipes = [];
let maxTimePass = 0;
let maxTimePassSec = 0;
let timePass = 0;
let reproductionCnt = 0;

let gameState = 0;
let classifier;

function setup() {
  createCanvas(640, 240);
  ml5.tf.setBackend("cpu");

  frameRate(60);
  classifier = ml5.neuralNetwork({
    inputs: 4,
    outputs: ["flap", "no flap"],
    task: "classification",
    neuroEvolution: true,
  });

  const modelInfo = {
    model: "model.json",
    metadata: "model_meta.json",
    weights: "model.weights.bin",
  };
  classifier.load(modelInfo, modelLoadedCallback);
}

function modelLoadedCallback() {
  console.log("modelLoaded ready");
  for (let i = 0; i < 10; i++) {
    birds[i] = new Bird(classifier);
  }
  pipes.push(new Pipe());
  gameState = 1;
}

function draw() {
  background(255);

  if (gameState == 1) {
    for (let i = pipes.length - 1; i >= 0; i--) {
      pipes[i].update();
      pipes[i].show();
      if (pipes[i].offscreen()) {
        pipes.splice(i, 1);
      }
    }

    for (let bird of birds) {
      if (bird.alive) {
        for (let pipe of pipes) {
          if (pipe.collides(bird)) {
            //noLoop();
            bird.alive = false;
          }
        }
        bird.think(pipes);
        bird.update();
        bird.show();
      }
    }

    if (frameCount % 100 == 0) {
      pipes.push(new Pipe());
    }

    if (allBirdsDead()) {
      normalizeFitness();
      reproduction();
      resetPipes();
      if (timePass > maxTimePass) {
        maxTimePass = timePass;
        maxTimePassSec = floor(timePass / 60);
      }

      timePass = 0;
      reproductionCnt += 1;
    } else {
      timePass += 1;
    }

    let timePassSec = floor(timePass / 60);

    push();
    fill(111);
    noStroke();
    textSize(16);
    text(
      `Round: ${reproductionCnt} MaxAlive: ${maxTimePassSec} Time: ${timePassSec}`,
      width - 240,
      20
    );
    pop();
  }
}

function allBirdsDead() {
  for (let bird of birds) {
    if (bird.alive) {
      return false;
    }
  }
  return true;
}

function reproduction() {
  let nextBirds = [];
  for (let i = 0; i < birds.length; i++) {
    let parentA = weightedSelection();
    let parentB = weightedSelection();
    let child = parentA.crossover(parentB);
    child.mutate(0.01);
    nextBirds[i] = new Bird(child);
  }
  birds = nextBirds;
}

function normalizeFitness() {
  let sum = 0;
  for (let bird of birds) {
    sum += bird.fitness;
  }
  for (let bird of birds) {
    bird.fitness = bird.fitness / sum;
  }
}

function weightedSelection() {
  let index = 0;
  let start = random(1);
  while (start > 0) {
    start = start - birds[index].fitness;
    index++;
  }
  index--;
  return birds[index].brain;
}

function resetPipes() {
  pipes.splice(0, pipes.length - 1);
}
