// Renders the landing terminal story as an animated SVG: node scripts/demo-svg.ts <out>

import { writeFileSync } from 'node:fs'
import { frames, type Line } from '../src/components/terminal/story.ts'

const out = process.argv[2]
if (!out) {
    console.error('usage: node scripts/demo-svg.ts <out.svg>')
    process.exit(2)
}

const ROWS = 18
const FONT = 14
const LINE = 22
const WIDTH = 760
const BAR = 34
const PAD_X = 18
const PAD_Y = 14
const BODY_TOP = BAR + PAD_Y
const HEIGHT = BODY_TOP + ROWS * LINE + PAD_Y

const COLORS: Record<string, string> = {
    green: '#89b29e',
    cyan: '#7fc8c8',
    yellow: '#e0af68',
    gray: '#9a949a',
    text: '#e0e0e0',
    dir: '#c5d4c8',
}

const escape = (s: string) => s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')

function tspans(line: Line): string {
    const segs = line.segs
        .filter(([, text]) => text)
        .map(([color, text, bold]) => {
            const weight = bold ? ' font-weight="700"' : ''
            return `<tspan fill="${COLORS[color]}"${weight}>${escape(text)}</tspan>`
        })
    if (line.cursor) segs.push(`<tspan class="c" fill="${COLORS.dir}">█</tspan>`)
    return segs.join('')
}

interface Shown {
    markup: string
    start: number
    end: number
}

// Each distinct content of a row, or of the title, and when it is on screen.
const shown: Shown[] = []
const open = new Map<string, Shown>()
const scroll: [number, number][] = []
let t = 0

const show = (slot: string, markup: string) => {
    const current = open.get(slot)
    if (current?.markup === markup) return
    if (current) current.end = t
    const next = { markup, start: t, end: Infinity }
    open.set(slot, next)
    shown.push(next)
}

for (const frame of frames()) {
    const seen = new Set(['title'])
    show('title', `<text class="t" x="${WIDTH / 2}" y="${BAR / 2}">lancher · ${escape(frame.dir)}</text>`)
    frame.lines.forEach((line, row) => {
        const slot = `row${row}`
        seen.add(slot)
        const y = BODY_TOP + row * LINE + LINE / 2
        show(slot, `<text x="${PAD_X}" y="${y}">${tspans(line)}</text>`)
    })
    for (const [slot, item] of open) {
        if (!seen.has(slot)) {
            item.end = t
            open.delete(slot)
        }
    }
    const offset = Math.max(0, frame.lines.length - ROWS) * LINE
    if (scroll.at(-1)?.[1] !== offset) scroll.push([t, offset])
    t += frame.ms
}

const total = t
const pct = (ms: number) => `${+((ms / total) * 100).toFixed(3)}%`

const keyframes: string[] = []
const items = shown.map((item, i) => {
    const steps = [`${pct(item.start)}{opacity:1}`]
    if (item.start > 0) steps.unshift('0%{opacity:0}')
    if (item.end < total) steps.push(`${pct(item.end)}{opacity:0}`)
    keyframes.push(`@keyframes k${i}{${steps.join('')}}`)
    const final = item.end === Infinity ? ' f' : ''
    return item.markup.replace('<text', `<text class="l${final}" style="animation-name:k${i}"`)
})

const scrollSteps = scroll.map(([at, y]) => `${pct(at)}{transform:translateY(${y ? -y : 0}px)}`).join('')
const finalOffset = scroll.at(-1)![1]

const title = items.filter((item) => item.includes('class="t"'))
const body = items.filter((item) => !item.includes('class="t"'))

const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${WIDTH}" height="${HEIGHT}" viewBox="0 0 ${WIDTH} ${HEIGHT}" role="img" aria-label="lancher new creates a project from a template, and lancher add saves it back as a new template">
<style>
text{font-family:"Fira Code",ui-monospace,SFMono-Regular,Menlo,Consolas,"DejaVu Sans Mono",monospace;font-size:${FONT}px;dominant-baseline:central;white-space:pre}
.t{font-size:12px;fill:${COLORS.gray};text-anchor:middle}
.l{opacity:0;animation-duration:${total}ms;animation-timing-function:step-end;animation-iteration-count:infinite}
.c{animation:blink 1s step-end infinite}
#s{animation:scroll ${total}ms step-end infinite}
@keyframes blink{50%{opacity:0}}
@keyframes scroll{${scrollSteps}}
${keyframes.join('\n')}
@media (prefers-reduced-motion:reduce){.l,.c,#s{animation:none}.f{opacity:1}#s{transform:translateY(-${finalOffset}px)}}
</style>
<clipPath id="w"><rect width="${WIDTH}" height="${HEIGHT}" rx="10"/></clipPath>
<clipPath id="b"><rect y="${BODY_TOP}" width="${WIDTH}" height="${ROWS * LINE}"/></clipPath>
<g clip-path="url(#w)">
<rect width="${WIDTH}" height="${HEIGHT}" fill="#141214"/>
<rect width="${WIDTH}" height="${BAR}" fill="#1b181b"/>
<rect y="${BAR - 1}" width="${WIDTH}" height="1" fill="#2e2a2e"/>
${[0, 1, 2].map((i) => `<circle cx="${18 + i * 20}" cy="${BAR / 2}" r="6" fill="#3f3a3f"/>`).join('')}
${title.join('\n')}
<g clip-path="url(#b)"><g id="s">
${body.join('\n')}
</g></g>
</g>
<rect x="0.5" y="0.5" width="${WIDTH - 1}" height="${HEIGHT - 1}" rx="10" fill="none" stroke="#3a353a"/>
</svg>
`

writeFileSync(out, svg)
