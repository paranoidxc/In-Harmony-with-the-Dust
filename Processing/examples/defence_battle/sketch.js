let ft_size, bow_len, arw_h, arw_w, zg_size;
let arrows = [];
let zergs = [];
let w, h, r, g, b, a, n;

function preload() {
  arrowS = loadSound("arrow.wav");
  hitS = loadSound("hitflesh.ogg");
  boomS = loadSound("boom.wav");
  roarS = loadSound("roar.wav");
}

function setup() {
  createCanvas(480, 480);
  w = width;
  h = height;
  ft_size = h / 20;
  bow_len = 1.05 * ft_size;
  arw_h = 0.8 * bow_len;
  arw_w = 0.2 * bow_len;
  zg_size = ft_size / 2;
  fort = new Fort(w / 2, h / 2);
  info = new Game(w / 30, w / 15);
  r = 255;
  g = 255;
  b = 255;
  a = 255;
  n = 10; // 虫子数量随机参数
  strokeWeight(3);
}

function draw() {
  background(0);
  if (info.status == 1) {
    fort.update();
    fort.show();
    if (arrows.length > 0) {
      for (let i = 0; i < arrows.length; i++) {
        if (arrows[i].timer < 1) {
          arrows.splice(i, 1);
        }
      }
      for (let i = 0; i < arrows.length; i++) {
        arrows[i].update();
        arrows[i].show();
      }
    }
    zergClean();
    zergHatch();
    for (let i = 0; i < zergs.length; i++) {
      push();
      zergs[i].update();
      zergs[i].show();
      pop();
    }
  }
  info.show();
}

// 点击鼠标开启或重启游戏，以及在游戏中发射箭矢
function mouseClicked() {
  if (info.status == 0) {
    info.status = 1;
    t = frameCount;
  } else if (info.status == 1 && info.hp <= 0) {
    // noLoop();
    info.status = 2;
  } else if (info.status == 2) {
    background(0);
    setup();
    arrows = [];
    zergs = [];
    info.status = 1;
    t = frameCount;
  } else {
    arrows.push(new Arrow(-bow_len, fort.angle));
    fort.color == "white" ? (fort.color = "yellow") : (fort.color = "white");
    arrowS.play();
  }
}

// 按空格键释放辐射圈
function keyPressed() {
  if (info.status == 1 && keyCode === 32 && fort.radiate == 0) {
    fort.radiate = 1;
    roarS.play();
  }
  if (key == "s" && zergs.length > 0) {
    for (let i = 0; i < zergs.length; i++) {
      zergs[i].alive = false;
    }
    hitS.play();
  }
}

// 检查箭矢是否射中虫子，参数x,y为箭矢当前坐标
function isHit(x, y) {
  for (let i = 0; i < zergs.length; i++) {
    if (dist(x, y, zergs[i].x, zergs[i].y) <= zergs[i].size / 2) {
      if (zergs[i].alive) {
        hitS.play();
      }
      zergs[i].alive = false;
      return true;
    }
  }
}

// 在四边界随机产生虫子
function zergHatch() {
  if ((frameCount - t - 1) % 120 == 0) {
    for (let i = 0; i < random(n); i++) {
      zergs.push(new Zerg(0, random(h), zg_size));
    }
    for (let i = 0; i < random(n); i++) {
      zergs.push(new Zerg(w, random(h), zg_size));
    }
    for (let i = 0; i < random(n); i++) {
      zergs.push(new Zerg(random(w), 0, zg_size));
    }
    for (let i = 0; i < random(n); i++) {
      zergs.push(new Zerg(random(w), h, zg_size));
    }
  }
}

// 清除虫子尸体
function zergClean() {
  if (zergs.length > 0) {
    for (let i = 0; i < zergs.length; i++) {
      if (zergs[i].timer < 1) {
        zergs.splice(i, 1);
      }
    }
  }
}
