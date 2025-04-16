<script setup>
	import {ref, watch, onMounted} from 'vue'
	import RollButton from './components/RollButton.vue'
	import socket from './lib/socket'
	import Die from './components/Die.vue'

	socket.Emitter.addEventListener("message", () => { })

	const state = ref(null)
	const dieRefs = ref([])
	const rolling = ref(false)
	const username = ref("")

	function ConnectWithId() {
		const url = new URL(window.location.href);
		socket.Id = url.searchParams.get('id') ?? "none"
		if (socket.Id === "none") {
			return
		}
		socket.Connect(socket.Id)
	}
	function ConnectWithUsername() {
		socket.Connect("none", username.value)
	}

	socket.Emitter.addEventListener("message", (event) => {
		const url = new URL(window.location.href);
		const callId = event.detail.callId;
		switch (callId) {
			case 0:
				if (event.detail.message.error !== undefined) {
					socket.Id = "none";
					url.searchParams.delete('id');
					window.history.pushState({}, '', url);
					state.value = null;
					console.log("Close")

				} else {
					socket.Id = event.detail.message.id
					url.searchParams.set('id', socket.Id);
					window.history.pushState({}, '', url);
				}
				break
			case 1:
				state.value = event.detail.message
				break
			case 2:
				state.value = event.detail.message
				rolling.value = true;
				dieRefs.value.forEach((dieComponent) => {
					dieComponent.roll();
				});
				setTimeout(() => {
					rolling.value = false;
				}, 750)
				break
		}
	});



	function onRoll() {
		rolling.value = true
	}

	onMounted(ConnectWithId);

	watch(
		state,
		(newVal, _) => {
			console.log('State changed:', newVal)
		},
		{deep: true}
	)
</script>

<template>
	<div v-if="state == null" class="grow flex flex-col gap-4 justify-center items-center">
		<input type="text" placeholder="Username" v-model="username" class="input" />
		<button class="btn" @click="ConnectWithUsername">Connect</button>
	</div>
	<div v-else class="grow flex flex-col">
		<p>{{state.players[socket.Id]}}</p>
		<div class="w-full flex flex-row gap-8 justify-center p-8">
			<Die v-for="die in state.dice" :key="die.index" :die="die" ref="dieRefs" />
		</div>
		<RollButton @roll="onRoll" :rolling="rolling" />

	</div>
</template>
