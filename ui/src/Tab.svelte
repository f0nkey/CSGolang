<script>
    import { createEventDispatcher } from 'svelte'
    export let chosen;
    export let title;
    $: currentlyChosen = chosen === title;

    const dispatch = createEventDispatcher();

   function sendEvent(){
       console.log("change it");
       chosen = title;
       dispatch("selection",title);
   }


</script>


    {#if currentlyChosen} <div class="tab tab-active"><span>{title}</span></div> {/if}
    {#if !currentlyChosen} <div on:click={sendEvent} class="tab"><span>{title}</span></div> {/if}


<style>
	.tab {
	  border-top-left-radius: 25px;
	  border-bottom-left-radius: 25px;
	  text-align:center;
      display: block;
      background-color: white;
      color: black;
      padding: 22px 0 22px 0;
      margin-bottom:25px;


      width: 100%;
      outline: none;
      cursor: pointer;
      transition: 0.3s;
    }
    .tab-active:before {
          background-color: red;
          color: red;
          content: "5";
          display: block;
          position: relative;
          float: right;
          left: 20px;
          margin-right: auto;
          height: 65px;
          top: -22px;
          width: 21px;
       }
    .tab-active {
    	  border-top-left-radius: 25px;
          border-bottom-left-radius: 25px;
          text-align:center;
          display: block;
          background-color: red;
          color: white;
          padding: 22px 0 22px 0;
          margin-bottom:25px;

          width: 100%;
          outline: none;
          cursor: pointer;
          transition: 0.3s;
          z-index: 3000;
          box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);

     }
    .tab-active span{
         animation-name: text-fly;
         animation-duration: 0.25s;
         animation-timing-function: ease-in-out;
    }

    @keyframes text-fly {
         from {filter: blur(0px);color:red;}
         to {filter: blur(0px);color:white;}
    }


</style>