<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import ErrorAlert from '$lib/components/ErrorAlert.svelte';
  import { zPodSchema, zPodServerSchema, zPodServerWithPodsSchema } from '$lib/models.js';
  import { pb, simpleSend, updateFromEvent } from '$lib/pocketbase.js';
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
    sistema para ter o seu limite ajustado.
  </p>

  <ErrorAlert errors={allErrors} />

  <div class="prose">
    <ul>
      {#each data.podServers as srv (srv.id)}
        {@const pods = data.pods.filter((p) => p.server === srv.id)}
        <li>
          <span><b>{srv.name}:</b> {srv.address}</span>
          <ul>
            {#each pods as pod (pod.id)}
              <li>
                <div class="flex justify-between">
                  <span>
                    <a href="pods/{pod.id}"><b>{pod.name} - {pod.image}:</b></a>
                    {pod.status}
                    {#if pod.running}
                      <span class="badge badge-success">running</span>
                    {:else}
                      <span class="badge badge-ghost">stopped</span>
                    {/if}
                  </span>

                  <div>
                    <button
                      type="button"
                      class="btn btn-info btn-sm"
                      disabled={anyLoading}
                      on:click={() => sendInspect(pod.id)}
                    >
                      inspect
                    </button>

                    <button
                      type="button"
                      class="btn btn-primary btn-sm"
                      disabled={anyLoading}
                      on:click={() => sendStart(pod.id)}
                    >
                      start
                    </button>

                    <button
                      type="button"
                      class="btn btn-accent btn-sm"
                      disabled={anyLoading}
                      on:click={() => sendStop(pod.id)}
                    >
                      stop
                    </button>

                    <button
                      type="button"
                      class="btn btn-warning btn-sm"
                      disabled={anyLoading}
                      on:click={() => sendKill(pod.id)}
                    >
                      kill
                    </button>
                  </div>
                </div>
              </li>
            {/each}
          </ul>
        </li>
      {/each}
    </ul>
  </div>
</BasicCard>
