
<script>
import { onMount } from 'svelte'
import Panel from './Panel.svelte';
import Tab from './Tab.svelte'
import TabLogo from './TabLogo.svelte'
import TabsList from './TabsList.svelte'
import Switch from './Switch.svelte'



let config;

getConfig();

$: {config;postConfig()}

setTimeout(function(){
    console.log("this is skeleton", config.toggles.skeleton);
    config.toggles.skeleton = false;
    config.toggles.skeleton = true;


},3000);


function getConfig(){

    var request = new XMLHttpRequest();
    request.open('GET', 'http://localhost:9991', false);  // `false` makes the request synchronous
    request.send(null);

    if (request.status === 200) {
        console.log("success grab");
        console.log(request.responseText);
        config = JSON.parse(request.responseText);
    } else {
        console.log("FAILURE");
    }
}

function postConfig(){
    console.log("posting ..");
    fetch("http://localhost:9991", {
         method: 'POST',
         headers: {},
         body: JSON.stringify(config)
    })
    .then(function (data) {
         console.log('POST success: ', data);
    })
    .catch(function (error) {
        console.log('POST failure: ', error);
    });
}

    let chosen = "General";

    function setSelectedTab(event){
        chosen = event.detail;
    }
</script>

<!--<div id="loading-text">Waiting for CSGolang.exe execution ..</div>!-->


<div class="grid" id="app">

    <TabsList>
        <TabLogo></TabLogo>
        <Tab on:selection={setSelectedTab} chosen={chosen} title="General"></Tab>
        <Tab on:selection={setSelectedTab} chosen={chosen} title="ESP"></Tab>
        <Tab on:selection={setSelectedTab} chosen={chosen} title="Misc"></Tab>
    </TabsList>

    <div class="config-display">
        {#if (chosen == "General")}
            <Switch bind:toggler={config.toggles.skeleton}></Switch> //todo: start aligning this working config toggles
            <h1>Yee haw</h1>
        {/if}

        {#if (chosen == "ESP")}
            <Panel title="Skeleton" bind:toggler={config.toggles.skeleton}> </Panel>
        {/if}

        {#if (chosen == "Misc")}
            <h1>Yee haw3</h1>
        {/if}

    </div>
</div>
<style>
    #app {
        display: grid;
        visibility: visible;
    }
	.grid {
	    width:100%;
	    height:100%;
	    display:grid;
        grid-template-columns: 25% 75%;
        grid-template-rows: auto;
	}


    .config-display {
        color:white;
        display:grid;
        grid-template-columns: 33% 33% 33%;
        background-color: red;
        width:100%;
        height:100%;
        box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    }
    .config-display:before {
        background-color: red;
        color: red;
        content: "5";
        display: block;
        position: absolute;
        height: 100%;
        width: 21px;
    }

    #loading-text {
        display:flex;
        justify-content:center;
        align-items:center;
        font-size: 40px;
    }
</style>
