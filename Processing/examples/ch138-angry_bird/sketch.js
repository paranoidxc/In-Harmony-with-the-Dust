const { Engine, World, Bodies, Mouse, MouseConstraint, Constraint } = Matter;

let ground;
let box;
const boxes = [];
let gird;
let world, engine;
let mConstraint;
let slingshot;

let dotImg;
let boxImg;
let bkgImg;

function preload() {
  dotImg = loadImage("images/dot.png");
  boxImg = loadImage("images/equals.png");
  bkgImg = loadImage("images/skyBackground.png");
}

function setup() {
  createCanvas(800, 600);
  engine = Matter.Engine.create();
  world = engine.world;

  ground = new Ground(width / 2, height - 10, width, 20);
  for (let i = 0; i < 5; i++) {
    boxes[i] = new Box((4 * width) / 5, 300 - i * 75, 50, 75);
  }
  bird = new Bird(150, (2 * height) / 3, 16);
  slingshot = new SlingShot(150, (2 * height) / 3, bird.body);

  const mouse = Mouse.create(canvas.elt);
  const options = {
    mouse: mouse,
  };
  //mouse.pixelRatio = pixelDensity();
  mConstraint = MouseConstraint.create(engine, options);
  World.add(world, mConstraint);
}

function keyPressed() {
  if (key == " ") {
    World.remove(world, bird.body);
    bird = new Bird(150, (2 * height) / 3, 16);
    slingshot.attach(bird.body);
  }
}

function mouseReleased() {
  setTimeout(() => {
    slingshot.fly();
  }, 10);
}

function draw() {
  background(0);
  Matter.Engine.update(engine);
  for (let box of boxes) {
    box.show();
  }
  slingshot.show();
  bird.show();
  ground.show();
}
