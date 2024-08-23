<script lang="ts">
  import BasicCard from '$lib/components/BasicCard.svelte';
  import { zEndDeviceSchema } from '$lib/models.js';
  import { pb, updateFromEvent } from '$lib/pocketbase.js';
  import { onMount } from 'svelte';
  import type { z } from 'zod';

  type EndDeviceModel = z.infer<typeof zEndDeviceSchema>;

  export let data;

  // --------------------------------------------------------------------------

  let unsub: () => void;
  const sub = async () => {
    unsub = await pb
      .collection('endDevices')
      .subscribe<EndDeviceModel>(
        '*',
        (e) => (data.devices = updateFromEvent(e, zEndDeviceSchema, data.devices)),
      );
  };

  // --------------------------------------------------------------------------

  onMount(() => {
    sub();

    return () => {
      if (unsub) unsub();
    };
  });
</script>

<BasicCard>
  <svelte:fragment slot="title">
    {#if data.user}
      {data.user.name} - {data.user.mat} - {data.user.curso}
    {:else}
      ...
    {/if}
  </svelte:fragment>
  {#if data.user}
    <h5 class="text-xs italic text-slate-400">@{data.user.username} - {data.user.email}</h5>

    <div class="h-2"></div>

    <p>
      Ao trocar email ou senha, um email automatico vai ser enviado para {data.user.email} com instrucoes.
    </p>

    <div class="flex gap-2">
      <button type="button" class="btn btn-warning btn-sm">trocar email</button>
      <button type="button" class="btn btn-warning btn-sm">trocar senha</button>
    </div>

    <div class="divider"></div>

    <h4 class="card-title">Seus dispositivos</h4>
    <ul class="mx-5">
      {#each data.devices as dev (dev.id)}
        <li class="border-b py-5 last:border-b-0">
          <div>
            <span>
              <b>{dev.deviceData?.system.browser.name} {dev.deviceData?.system.browser.version}</b>
              -
              <b>{dev.deviceData?.system.platform}</b>
            </span>
            <span class="text-sm">{dev.locationData?.city}</span>
          </div>
        </li>
      {/each}
    </ul>
  {:else}
    ...
  {/if}
</BasicCard>
