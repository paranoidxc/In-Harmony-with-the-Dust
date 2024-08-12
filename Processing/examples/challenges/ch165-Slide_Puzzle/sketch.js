let source;
let tiles = [];
let cols = 4;
let rows = 4;
let w, h;
let board = [];
let blank = 0;
let simpleShuffleNum = 4;

function preload() {
  //source = loadImage("choochoobot.png");
  source = loadImage("demo.jpeg");
}

function setup() {
  createCanvas(400, 400);
  w = width / cols;
  h = height / rows;

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < rows; j++) {
      let x = i * w;
      let y = j * h;
      let img = createImage(w, h);
      img.copy(source, x, y, w, h, 0, 0, w, h);
      //img.loadPixels();
      //img.updatePixels();
      let index = i + j * cols;
      board.push(index);
      let tile = new Tile(index, img);
      tiles.push(tile);
    }
  }

  tiles.pop();
  board.pop();
  board.push(-1);

  simpleShuffle(board);
}

function swap(i, j, arr) {
  let tmp = arr[i];
  arr[i] = arr[j];
  arr[j] = tmp;
}

function randomMove(arr) {
  let r1 = floor(random(0, cols));
  let r2 = floor(random(0, rows));
  move(r1, r2, arr);
}

function simpleShuffle(arr) {
  for (let i = 0; i < simpleShuffleNum; i++) {
    randomMove(arr);
  }
}

function mousePressed() {
  let i = floor(mouseX / w);
  let j = floor(mouseY / h);
  move(i, j, board);
}

function isSolved() {
  for (let i = 0; i < board.length - 1; i++) {
    if (board[i] !== tiles[i].index) {
      return false;
    }
  }
  return true;
}

function move(i, j, arr) {
  let blank = findBlank();
  let blankCol = blank % cols;
  let blankRow = floor(blank / rows);

  if (isNeighbor(i, j, blankCol, blankRow)) {
    swap(blank, i + j * cols, arr);
  }
}

function isNeighbor(i, j, x, y) {
  if (i != x && j != y) {
    return false;
  }
  if (abs(i - x) == 1 || abs(j - y) == 1) {
    return true;
  }
  return false;
}

function findBlank() {
  for (let i = 0; i < board.length; i++) {
    if (board[i] == -1) {
      return i;
    }
  }
}

function draw() {
  background(0);
  //randomMove(board);
  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < rows; j++) {
      let index = i + j * cols;
      let x = i * w;
      let y = j * h;
      let tileIndex = board[index];
      if (tileIndex > -1) {
        let imgx = tiles[tileIndex].img;
        image(imgx, x, y, w, h);
      }
    }
  }

  for (let i = 0; i < tiles.length; i++) {
    for (let j = 0; j < rows; j++) {
      let x = i * w;
      let y = j * h;
      strokeWeight(2);
      noFill();
      rect(x, y, w, h);
    }
  }

  if (isSolved()) {
    console.log("SOLVED");
  }
}
