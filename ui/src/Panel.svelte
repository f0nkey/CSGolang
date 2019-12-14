<script>
    import Switch from './Switch.svelte'
    export let title;
    export let toggler;

</script>


<div>
    <div style="position: relative; width: 0; height: 0;margin:0; margin-left: calc(100% - 70px); top: 4px;">
        <Switch bind:checked={toggler}></Switch>
    </div>
    <h1>{title}</h1>


    <hr>
    <h1>what</h1>
    <slot></slot>

</div>

<style>
    div {
        margin: 20px;
        color:black;
        background-color:white;
        border-radius: 25px;
        box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    }
    h1 {
        padding-left: 10px;
        font-size: 16px;
    }
    hr {
        width:100%;
    }

    input[type=checkbox]{
    	height: 0;
    	width: 0;
    	visibility: hidden;
    }

    label {
    	cursor: pointer;
    	text-indent: -9999px;
    	width: 50px;
    	height: 25px;
    	background: grey;
    	display: block;
    	border-radius: 100px;
    	position: relative;
    	top:-23px;
    	left: 5px;

    }

    label:after {
    	content: '';
    	position: absolute;
    	left: 1px;
    	top:1px;
    	width: 23px;
    	height: 23px;
    	background: #fff;
    	border-radius: 24px;
    	transition: 0.3s;
    }

    input:checked + label {
    	background: darkred;
    }

    input:checked + label:after {
    	left: calc(100% - 3px);
    	transform: translateX(-100%);
    }

    label:active:after {
    	width: 40px;
    }
</style>
