
<template>
  <div id="app">
    <div id="FamilyChart" class="f3" style="width:100%;height:900px;margin:auto;background-color:rgb(33,33,33);color:#fff;"></div>
  </div>
</template>

<script>
import ft from './data.json';
import f3 from "./family-chart/dist/family-chart.js";
import './family-chart.css';  // create file 'family-chart.css' in same directory, copy/paste css from examples/create-tree

export default {
  mounted() {
  // fetch('https://raw.githubusercontent.com/giatoc/family-tree/refs/heads/master/src/components/data.json')
  // .then(res => res.json())
  // .then(data => create(data))
  // .catch(err => console.error(err))

    create(data())

    function create(data) {
      const cont = document.querySelector("div#FamilyChart")  // make sure to create div with id FamilyChart
      const store = f3.createStore({
        data,
        node_separation: 250,
        level_separation: 150
      })
      const svg = f3.createSvg(cont)
      const Card = f3.elements.Card({
        store,
        svg,
        card_dim: {w:220,h:70,text_x:75,text_y:15,img_w:60,img_h:60,img_x:5,img_y:5},
        card_display: [d => `${d.data['first name'] || ''} ${d.data['last name'] || ''}`, d => `${d.data['birthday'] || ''}`],
        mini_tree: true,
        link_break: false
      })

      store.setOnUpdate(props => f3.view(store.getTree(), svg, Card, props || {}))
      store.updateTree({initial: true})
    }

    function data() {
      return ft
    }
  }
};
</script>
<style></style>