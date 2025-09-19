<script lang="ts">
  import { getSourceServiceClient } from "#/connectrpc";
  import { toDate } from "#/connectrpc/js_date";
  import { Popover, Combobox } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";
  import { resource, watch } from "runed";
  import { theme } from "#/store/theme";
  import cronstrue from "cronstrue";

  interface Props {
    value: string[];
  }

  let scheduleInputValue = $state("");
  let cronExpressionError = $derived.by(() => {
    if (!scheduleInputValue.trim()) {
      return "";
    }
    try {
      cronstrue.toString(scheduleInputValue, {
        throwExceptionOnParseError: true,
      });
      return "";
    } catch (e) {
      return e instanceof Error ? e.message : `${e}`;
    }
  });

  const scheduleInputNextRunResource = resource(
    () => scheduleInputValue,
    async (pattern, _, { signal }) => {
      if (!pattern.trim()) {
        return;
      }
      return getNextRun(pattern, signal);
    },
    {
      debounce: 300,
    },
  );

  watch(
    () => scheduleInputNextRunResource.error,
    (err) => {
      if (err) {
        cronExpressionError = err instanceof Error ? err.message : `${err}`;
      }
    },
  );

  async function getNextRun(pattern: string, signal?: AbortSignal) {
    cronstrue.toString(pattern);
    const client = await getSourceServiceClient({ signal });
    return client.getCronNextTime({
      cronExpression: pattern,
    });
  }

  let { value = $bindable([]) }: Props = $props();

  type ScheduleEntry = {
    pattern: string;
    nextRun?: Date;
    error?: string;
    zone?: string;
  };

  const schedules: Promise<ScheduleEntry[]> = $derived(
    Promise.allSettled(value.map((pattern) => getNextRun(pattern))).then(
      (results) =>
        results.map((res, idx) =>
          res.status === "fulfilled"
            ? {
                pattern: value[idx],
                nextRun: toDate(res.value.nextTime!),
                zone: res.value.zone || "UTC",
              }
            : {
                pattern: value[idx],
                error:
                  res.reason instanceof Error
                    ? res.reason.message
                    : `${res.reason}`,
              },
        ),
    ),
  );

  type PatternOption = {
    value: string;
    label: string;
  };

  function toPatternOption(value: string): PatternOption {
    let label;
    try {
      label = cronstrue.toString(value, {
        throwExceptionOnParseError: true,
        verbose: true,
        use24HourTimeFormat: true,
      });
    } catch (e) {
      label = `Invalid pattern: ${e instanceof Error ? e.message : e}`;
    }
    return { value, label };
  }

  function genDropdownDays(): PatternOption[] {
    const hours = Array.from({ length: 24 }).keys();
    const dayOfWeeks = ["SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"];
    const patterns: PatternOption[] = [];

    for (const hour of hours) {
      for (const day of dayOfWeeks) {
        patterns.push(toPatternOption(`0 ${hour} * * ${day}`));
      }
      patterns.push(toPatternOption(`0 ${hour} * * *`));
    }
    return patterns;
  }

  const patternDropdownList: PatternOption[] = [
    toPatternOption("* * * * *"),
    toPatternOption("*/5 * * * *"),
    toPatternOption("*/15 * * * *"),
    toPatternOption("*/30 * * * *"),
    toPatternOption("0 * * * *"),
    toPatternOption("0 0 * * MON,TUE,WED,THU,FRI"),
  ].concat(genDropdownDays());

  const filteredPatternDropdownList = $derived.by(() => {
    const text = scheduleInputValue.trim().toLowerCase();
    if (!text) {
      return patternDropdownList;
    }
    const filtered = patternDropdownList.filter(
      (opt) =>
        opt.value.toLowerCase().includes(text) ||
        opt.label.toLowerCase().includes(text),
    );
    return filtered;
  });

  let showDropdown = $state(false);

  let inputField: HTMLInputElement | null = $state(null);

  let isMobile = $state(window.innerWidth < 640);
</script>

<svelte:window onresize={() => (isMobile = window.innerWidth < 640)} />

<Combobox.Root
  type="single"
  inputValue={scheduleInputValue}
  onValueChange={(value) => (scheduleInputValue = value.trim())}
  open={showDropdown}
>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">
      <span>Schedule</span>
    </legend>
    <Combobox.Input
      placeholder="e.g. 0 0 * * FRI ({toPatternOption('0 0 * * FRI').label})"
      oninput={(e) => (scheduleInputValue = e.currentTarget.value)}
      onfocus={() => (showDropdown = true)}
      onblur={() => (showDropdown = false)}
      class="input w-full"
      bind:ref={inputField}
    ></Combobox.Input>
    {#if cronExpressionError}
      <div class="alert alert-error alert-soft">
        {cronExpressionError}
      </div>
    {:else if scheduleInputNextRunResource.loading}
      <div class="alert alert-warning alert-soft">
        <div class="loading loading-spinner"></div>
        <span>Validating cron exppression...</span>
      </div>
    {:else if scheduleInputNextRunResource.current}
      {@const data = scheduleInputNextRunResource.current}
      {@const intl = new Intl.DateTimeFormat().resolvedOptions()}
      {@const ts = toDate(data.nextTime!)}
      {@const nextRunLocal = ts.toLocaleString()}
      {@const nextRunServer = ts.toLocaleString(intl.locale, {
        timeZone: data.zone,
      })}
      <div class="alert alert-success alert-soft flex flex-col">
        <span>
          Next Run (local time): {nextRunLocal}
        </span>
        <span>Server Next Run: {nextRunServer}</span>
      </div>
    {:else}
      <p class="label">Cron expression pattern. Only 5 fields are supported.</p>
    {/if}
  </fieldset>
  <Combobox.Portal>
    <Combobox.Content
      data-theme={$theme}
      class="z-[9999]"
      side="top"
      sideOffset={10}
    >
      <Combobox.Viewport
        class="menu bg-base-200 rounded-sm shadow-[0_-0.5rem_1.5rem_rgba(0,0,0,0.1)]"
        style={{ width: inputField ? `${inputField.offsetWidth}px` : "auto" }}
      >
        <div
          class="flex flex-col overflow-auto"
          style={`width: ${inputField ? `${inputField.offsetWidth}px` : "auto"}; max-height: ${isMobile ? "12.5rem" : "20rem"}`}
        >
          {#each filteredPatternDropdownList as item, i (item.value)}
            {#if i > 0}
              <div class="divider py-0 my-0"></div>
            {/if}
            <Combobox.Item
              value={item.value}
              label={item.value}
              class="btn btn-ghost flex-col items-start justify-center gap-0 py-8"
            >
              <span>{item.value}</span>
              <span class="font-normal text-left">{item.label}</span>
            </Combobox.Item>
          {:else}
            <Combobox.Item
              value={scheduleInputValue}
              label={scheduleInputValue}
              class="btn btn-ghost flex-col items-start justify-center gap-0 py-8"
            >
              <span>{scheduleInputValue}</span>
              <span class="font-normal text-left">
                {toPatternOption(scheduleInputValue).label}
              </span>
            </Combobox.Item>
          {/each}
        </div>
      </Combobox.Viewport>
    </Combobox.Content>
  </Combobox.Portal>
</Combobox.Root>

{#snippet helpText()}
  <Popover.Root>
    <Popover.Trigger class="btn btn-square btn-ghost btn-xs">
      <IconInfo />
    </Popover.Trigger>
    <Popover.Portal>
      <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula"
      ></Popover.Content>
    </Popover.Portal>
  </Popover.Root>
{/snippet}
