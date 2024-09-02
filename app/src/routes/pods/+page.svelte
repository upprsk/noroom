<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { zPodSchema, zPodServerSchema, zPodServerWithPodsSchema } from '$lib/models';
  import { pb, simpleSend, updateFromEvent } from '$lib/pocketbase';
  import { onMount } from 'svelte';
  import { z } from 'zod';

  export let data;

  const {
    loading: loadingStart,
    errors: errorsStart,
    send: sendStart,
  } = simpleSend(pb, z.unknown(), (id: string) => `/api/noroom/pod/${id}/start`, {
    method: 'post',
  });

  const {
    loading: loadingStop,
    errors: errorsStop,
    send: sendStop,
  } = simpleSend(pb, z.unknown(), (id: string) => `/api/noroom/pod/${id}/stop`, {
    method: 'post',
  });

  const {
    loading: loadingKill,
    errors: errorsKill,
    send: sendKill,
  } = simpleSend(pb, z.unknown(), (id: string) => `/api/noroom/pod/${id}/kill`, {
    method: 'post',
  });

  const {
    loading: loadingInspect,
    errors: errorsInspect,
    send: sendInspect,
  } = simpleSend(pb, z.unknown(), (id: string) => `/api/noroom/pod/${id}/inspect`, {
    method: 'post',
  });

  $: anyLoading = $loadingStart || $loadingStop || $loadingKill || $loadingInspect;
  $: allErrors = [
    ...($errorsStart ?? []),
    ...($errorsStop ?? []),
    ...($errorsKill ?? []),
    ...($errorsInspect ?? []),
  ];

  // --------------------------------------------------------------------------

  type PodServer = z.infer<typeof zPodServerWithPodsSchema>;
  type Pod = z.infer<typeof zPodSchema>;

  let unsubPodServers: () => void;
  const subPodServers = async () => {
    unsubPodServers = await pb
      .collection('podServers')
      .subscribe<PodServer>(
        '*',
        (e) => (data.podServers = updateFromEvent(e, zPodServerSchema, data.podServers)),
      );
  };

  // --------------------------------------------------------------------------

  let unsubPods: () => void;
  const subPods = async () => {
    unsubPods = await pb
      .collection('pods')
      .subscribe<Pod>('*', (e) => (data.pods = updateFromEvent(e, zPodSchema, data.pods)));
  };

  // --------------------------------------------------------------------------

  onMount(() => {
    subPodServers();
    subPods();

    return () => {
      if (unsubPodServers) unsubPodServers();
      if (unsubPods) unsubPods();
    };
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">Servidores Pod</svelte:fragment>

  <div class="h-5"></div>

  <p>
    Os servidores Pod permitem criar maquinas virtuais na nuvem. Cada usuario tem direito por padrao
    a uma maquina virtual. Se precisar de mais maquinas, entre em contato com um administrador do
    sistema para ter o seu limite ajustado. {#if data.user.role === 'editor'}Voce nao possui limite.{:else}Seu
      limite atual e de <b>{data.user.maxPods}</b> maquinas.{/if}
  </p>

  <ErrorAlert errors={allErrors} />

  {#if data.user.pods.length < data.user.maxPods || data.user.role === 'editor'}
    <div class="flex w-full justify-end">
      <a href="new" class="btn btn-primary btn-sm">novo</a>
    </div>
  {/if}

  <div class="overflow-x-auto">
    <table class="table">
      <!-- head -->
      <thead>
        <tr>
          <th>Nome</th>
          <th>Imagem</th>
          <th>Servidor</th>
          <th>Status</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each data.pods as pod (pod.id)}
          {@const srv = data.podServers.find((s) => s.id === pod.server)}
          <tr>
            <td>
              {pod.name}
              <br />
              {pod.podId.substring(0, 16)}
            </td>
            <td>
              {pod.image}
            </td>
            <td>
              {srv?.name}
            </td>
            <td>
              <span class="badge" class:badge-ghost={!pod.running} class:badge-success={pod.running}
                >{pod.status}</span
              >
            </td>
            <th>
              <button
                class="btn btn-secondary btn-xs"
                disabled={anyLoading}
                on:click={() => sendInspect(pod.id)}>update</button
              >

              {#if !pod.running}
                <button
                  class="btn btn-primary btn-xs"
                  disabled={anyLoading}
                  on:click={() => sendStart(pod.id)}>start</button
                >
              {:else}
                <button
                  class="btn btn-warning btn-xs"
                  disabled={anyLoading}
                  on:click={() => sendStop(pod.id)}>stop</button
                >
                <button
                  class="btn btn-warning btn-xs"
                  disabled={anyLoading}
                  on:click={() => sendKill(pod.id)}>kill</button
                >
              {/if}

              <a href={pod.id} class="btn btn-ghost btn-xs">details</a>
            </th>
          </tr>
        {/each}
      </tbody>
      <!-- foot -->
      <tfoot>
        <tr>
          <th>Nome</th>
          <th>Imagem</th>
          <th>Servidor</th>
          <th>Status</th>
          <th></th>
        </tr>
      </tfoot>
    </table>
  </div>
</BasicCard>
