<script lang="ts">
  export let name: string;
  export let errors: string[] | undefined;
  export let value: string;
  export let constraints: { [key: string]: any } | undefined;

  export let placeholder = name;
</script>

<label class="form-control w-full max-w-xs">
  <div class="label">
    <span class="label-text"><slot /></span>
  </div>
  <input
    type="password"
    class="input input-bordered w-full max-w-xs"
    aria-invalid={errors ? 'true' : undefined}
    bind:value
    {name}
    {placeholder}
    {...constraints}
  />
  <div class="label">
    {#if errors}<span class="label-text-alt text-error">{errors}</span>{/if}
  </div>
</label>
