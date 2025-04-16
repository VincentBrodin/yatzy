<script setup>
	import {ref} from 'vue'
	import RollButton from './components/RollButton.vue'
	import socket from './lib/socket'
	import Die from './components/Die.vue'

	socket.Emitter.addEventListener("message", () => { })

	const dice = ref(Array.from({length: 5}, (_, i) => ({index: i, value: 1, selected: false})));
	const dieRefs = ref([])

	socket.Emitter.addEventListener("message", (event) => {
		const callId = event.detail.callId;
		if (callId !== 2) return;
		dice.value = event.detail.message.dice
	});


	function onRoll(newDice) {
		console.log('New dice values:', newDice)
		dice.value = newDice
		dieRefs.value.forEach((dieComponent) => {
			dieComponent.roll()
		})
	}
</script>

<template>
	<div class="w-full flex flex-row gap-8 justify-center p-8">
		<Die v-for="die in dice" :key="die.index" :die="die" ref="dieRefs" />
	</div>
	<RollButton @roll="onRoll" />
</template>
