var cells = [];

function setup() {
  createCanvas(700, 700);
  cells.push(new Cell());
  //cells.push(new Cell());
}

function draw() {
  background(0);
  for (var i = 0; i < cells.length; i++) {
    cells[i].move();
    cells[i].show();
  }
  if (frameCount % 10 == 0) {
    //mitosis();
  }
}

function mousePressed() {
  mitosis();
}

function mitosis() {
  for (var i = cells.length - 1; i >= 0; i--) {
    //if (cells[i].clicked(mouseX, mouseY)) {
    cells.push(cells[i].mitosis());
    cells.push(cells[i].mitosis());
    cells.splice(i, 1);
    //}
  }
}
