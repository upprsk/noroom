<script lang="ts">
  import { page } from '$app/stores';
  import BasicCard from '$lib/components/BasicCard.svelte';
  import { pb } from '$lib/pocketbase';
  import { ClientResponseError } from 'pocketbase';
  import { onMount } from 'svelte';

  let presenceStatus = '';
  let presenceError = '';

  const getLocation = () => {
    presenceStatus = '';
    presenceError = '';

    if (!navigator.geolocation) {
      presenceStatus = 'Localizacao nao suportada pelo navegador';
      return;
    }

    presenceStatus = 'Procurando...';
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const { latitude, longitude } = pos.coords;
        presenceStatus = `Localizacao encontrada: ${latitude}, ${longitude}`;

        sendPresence({ latitude, longitude });
      },
      (e) => {
        presenceStatus = `Erro em procurar localizacao: ${e}`;
      },
      { enableHighAccuracy: true },
    );
  };

  const sendPresence = async (position: { latitude: number; longitude: number }) => {
    const { getFingerprint } = await import('@thumbmarkjs/thumbmarkjs');
    const fingerprint = await getFingerprint();

    try {
      const res = await pb.send('/api/noroom/presence', {
        method: 'POST',
        body: {
          class: $page.params.id,
          fingerprint,
          position,
        },
      });

      console.log(res);
      presenceStatus = `Presenca atualizada para ${res.user} (${res.dist}m)`;
    } catch (e) {
      if (e instanceof ClientResponseError) {
        presenceError = `Erro ao aplicar presenca: ${e.message}`;
        console.log(e.data);
      } else {
        presenceError = `Erro ao aplicar presenca: ${e}`;
      }
    }
  };

  export let data;

  onMount(() => {
    getLocation();
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">Aplicar presenca para {data.klass.title}</svelte:fragment>

  {#if presenceStatus}
    <div role="alert" class="alert">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        class="h-6 w-6 shrink-0 stroke-info"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        ></path>
      </svg>
      <span>{presenceStatus}</span>
    </div>
  {/if}

  {#if presenceError}
    <div role="alert" class="alert alert-error">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-6 w-6 shrink-0 stroke-current"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <span>{presenceError}</span>
    </div>
  {/if}

  <div class="card-actions justify-end">
    <a href="/classes/{data.klass.id}" class="btn">voltar</a>
  </div>
</BasicCard>
