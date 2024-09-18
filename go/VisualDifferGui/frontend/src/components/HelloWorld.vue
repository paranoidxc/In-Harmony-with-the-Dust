<script setup>
import {reactive} from 'vue'
import {Greet, SelectOld, SelectOldFolder, SelectNew, CallCompare} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

function greet() {
  Greet(data.name).then(result => {
    data.resultText = result
  })
}

function selectOld() {
  SelectOldFolder().then(result => {
    data.old = result
  })
}

function selectNew() {
  SelectNew().then(result => {
    data.new = result
  })
}

function compare() {
  CallCompare(data.old, data.new).then(result => {
    data.result = result
    console.log(result)
  })
} 


</script>

<template>
  <main>
    <!--
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="greet">Greet</button>
    </div>
    -->

    <div id="input" class="input-box">
      <input id="old" v-model="data.old" autocomplete="off" disabled class="input" type="text"/>
      <button class="btn" @click="selectOld">LEFT</button>
    </div>

    <div id="input" class="input-box">
      <input id="new" v-model="data.new" autocomplete="off" disabled class="input" type="text"/>
      <button class="btn" @click="selectNew">RIGHT</button>
    </div>

    <div id="input" class="input-box">
      <button class="btn" @click="compare">compare</button>
    </div>

    <textarea v-model="data.result" class="textareaResult"></textarea>
  </main>
</template>

<style scoped>
.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box {
  margin: 10px;
}
.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  width: 60%;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.textareaResult {
  border: none;
  border-radius: 3px;
  outline: none;
  line-height: 20px;
  padding: 0 10px;
  height: 600px;
  width: 70%;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}
</style>
