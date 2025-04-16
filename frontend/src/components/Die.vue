<script setup>
	import {ref} from 'vue';
	import socket from '@/lib/socket';

	const emit = defineEmits(["select"]);

	const props = defineProps({
		die: {
			type: Object,
			required: true
		}
	});

	const preview = ref(1);
	const rolling = ref(false);

	function getRandomInt(max) {
		return Math.floor(Math.random() * max);
	}


	let rollInterval = null
	function roll() {
		if (props.die.selected) return
		rolling.value = true
		if (rollInterval) clearInterval(rollInterval)

		rollInterval = setInterval(() => {
			preview.value = getRandomInt(6) + 1
		}, 100)

		setTimeout(() => {
			clearInterval(rollInterval)
			rollInterval = null
			preview.value = props.die.value
			rolling.value = false
		}, 750)
	}

	function select() {
		console.log(`Selected ${props.die.index}`)
		props.die.selected = !props.die.selected
		const payload = {index: props.die.index, selected: props.die.selected};
		socket.Send(2, JSON.stringify(payload));

	}
	defineExpose({roll})
</script>

<template>
	<div class="die-container">
		<button @click="select">
			<img class="w-full h-auto transition-opacity" :class="{ 'opacity-50': rolling || props.die.selected  }"
				:src="`/img/${preview}.png`" alt="Dice Face">
		</button>
	</div>
</template>

<style scoped>
	.die-container {
		max-width: 100%;
		aspect-ratio: 1 / 1;
	}
</style>
