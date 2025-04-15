<script setup>
	import {ref} from 'vue';
	import socket from '@/lib/socket';

	const emit = defineEmits(["roll"]);

	const rolling = ref(false)

	socket.Emitter.addEventListener("message", (event) => {
		const callId = event.detail.callId;
		if (callId !== 1) return;

		const dice = [0, 0, 0, 0, 0]

		for (let i = 0; i < event.detail.message.byteLength / 4; i++) {
			dice[i] = Number(event.detail.message.getInt32(i * 4));
		}

		emit("roll", dice)

		setTimeout(() => {
			rolling.value = false;
		}, 750)
	});

	function roll() {
		rolling.value = true;
		socket.Send(1);
	}

</script>

<template>
	<button class="btn" :disabled="rolling" @click="roll">
		<template v-if="rolling">
			Rolling
		</template>
		<template v-else>
			Roll
		</template>
	</button>
</template>
