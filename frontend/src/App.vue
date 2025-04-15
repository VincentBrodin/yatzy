<script setup>
	import {ref} from 'vue'
	import RollButton from './components/RollButton.vue'
	import socket from './lib/socket'
	import Die from './components/Die.vue'

	socket.Emitter.addEventListener("message", () => { })

	const dice = ref([1, 1, 1, 1, 1])

	const dieRefs = ref([])

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
		<Die v-for="(die, index) in dice" :key="index" :target="die" ref="dieRefs" />
	</div>
	<RollButton @roll="onRoll" />
</template>
