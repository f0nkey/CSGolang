
<script>
import { onMount } from 'svelte'
import Panel from './Panel.svelte';
import Tab from './Tab.svelte'
import TabLogo from './TabLogo.svelte'
import TabsList from './TabsList.svelte'
import Switch from './Switch.svelte'
import Select from './Select.svelte'



let config;

getConfig();

$: {config;postConfig()}

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
    config.colorModes.skeleton = parseInt(config.colorModes.skeleton);
    config.colorModes.name = parseInt(config.colorModes.name);
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

<div class="grid" id="app">
    <TabsList>
        <TabLogo></TabLogo>
        <Tab on:selection={setSelectedTab} chosen={chosen} title="General"></Tab>
        <Tab on:selection={setSelectedTab} chosen={chosen} title="ESP"></Tab>
        <Tab on:selection={setSelectedTab} chosen={chosen} title="Misc"></Tab>
    </TabsList>

    <div class="config-display">
        {#if (chosen == "General")}
            <h1>Yee haw</h1>
        {/if}

        {#if (chosen == "ESP")}
            <Panel title="Skeleton" bind:toggler={config.toggles.skeleton}>
                 <h1>See Teammates:</h1> <Switch bind:checked={config.seeTeammates.skeleton}></Switch>
                 <h1>Color Mode:</h1> <Select bind:selected={config.colorModes.skeleton}></Select>
            </Panel>

            <Panel title="Name" bind:toggler={config.toggles.name}>
                 <h1>See Teammates:</h1> <Switch bind:checked={config.seeTeammates.name}></Switch>
                 <h1>Color Mode:</h1> <Select bind:selected={config.colorModes.name}></Select>
            </Panel>
        {/if}

        {#if (chosen == "Misc")}
            <Panel title="Bhop" bind:toggler={config.toggles.bhop}>
                 <h1>Toggle:</h1> <Switch bind:checked={config.toggles.bhop}></Switch>
            </Panel>
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
	h1 {
	    padding-left: 10px;
        font-size: 16px;
        display:inline-block;
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
