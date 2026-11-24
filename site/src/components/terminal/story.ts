// The landing terminal: create a project, customize it, save it as a template.
// Output follows lancher v0.4 (internal/cli, internal/ui).

/** Terminal colors, as CSS classes `term-<color>`. */
export type Color = 'green' | 'cyan' | 'yellow' | 'gray' | 'text' | 'dir'

/** A run of text: color, text, bold. */
export type Seg = [Color, string, boolean?]

export type Step =
    | { cmd: string; dir: string }
    | { ask: string; answer: string }
    | { spin: string; ms: number }
    | { out: Seg[]; ms?: number }

export interface Line {
    segs: Seg[]
    cursor?: boolean
}

const HOME = '~/code'
const PROJECT = '~/code/my-api'

export const END_DIR = PROJECT

export const story: Step[] = [
    { cmd: 'lancher new go-api my-api', dir: HOME },
    { out: [] },
    { out: [['cyan', 'Template Information:', true]] },
    {
        out: [
            ['gray', '    Name:'],
            ['text', ' Go API'],
        ],
    },
    {
        out: [
            ['gray', '    Description:'],
            ['text', ' HTTP service with chi'],
        ],
    },
    { out: [] },
    { out: [['cyan', 'Hooks found:', true]] },
    { out: [['text', '  1. go mod tidy']] },
    { ask: 'Execute hooks? (y/N):', answer: 'y' },
    { ask: 'Initialize git repository? (y/N):', answer: 'y' },
    { out: [] },
    {
        out: [
            ['cyan', 'Lancher files in the template:', true],
            ['text', ' .lancher.yaml'],
        ],
    },
    { ask: 'Copy them to the project? (y/N):', answer: 'y' },
    { spin: 'Creating project...', ms: 900 },
    { out: [['green', "✓ Project created from template 'go-api'"]] },
    {
        out: [
            ['yellow', '  Location:'],
            ['text', ' /home/me/code/my-api'],
        ],
    },
    { out: [['green', '✓ Git repository initialized']] },
    { out: [] },
    {
        out: [
            ['cyan', 'Executing hook 1/1:'],
            ['text', ' go mod tidy'],
        ],
    },
    { out: [['green', '✓ All hooks executed successfully']], ms: 1200 },
    { cmd: 'cd my-api', dir: HOME },
    { cmd: 'go get github.com/jackc/pgx/v5', dir: PROJECT },
    { out: [['gray', 'go: added github.com/jackc/pgx/v5 v5.7.2']], ms: 500 },
    { cmd: 'cp -r ~/snippets/db internal/db', dir: PROJECT },
    { cmd: 'lancher add . go-api-pgx', dir: PROJECT },
    { spin: 'Copying template...', ms: 1000 },
    { out: [['green', "✓ Template 'go-api-pgx' added"]] },
    {
        out: [
            ['yellow', '  Source:'],
            ['text', ' /home/me/code/my-api'],
        ],
    },
    {
        out: [
            ['yellow', '  Stored:'],
            ['text', ' /home/me/.local/share/lancher/templates/go-api-pgx'],
        ],
    },
]

export function prompt(dir: string): Seg[] {
    return [
        ['dir', `${dir} `],
        ['green', '$ '],
    ]
}

/** One screen of the player: show `lines`, then wait `ms`. */
export interface Frame {
    lines: Line[]
    dir: string
    ms: number
}

const FRAMES = '-\\|/'
const HOLD = 5000

const keyDelay = (i: number) => 38 + ((i * 37) % 45)

/** The story as the player shows it, frame by frame, ending on the hold. */
export function* frames(): Generator<Frame> {
    const lines: Line[] = []
    let dir = HOME

    const snap = (ms: number): Frame => ({
        lines: lines.map((line) => ({
            segs: line.segs.map((seg) => [...seg] as Seg),
            cursor: line.cursor,
        })),
        dir,
        ms,
    })

    const push = (segs: Seg[], cursor = false) => {
        for (const line of lines) line.cursor = false
        const line = { segs, cursor }
        lines.push(line)
        return line
    }

    yield snap(400)
    for (const step of story) {
        if ('cmd' in step) {
            dir = step.dir
            const typed: Seg = ['text', '']
            const line = push([...prompt(step.dir), typed], true)
            yield snap(500)
            for (let i = 0; i < step.cmd.length; i++) {
                typed[1] += step.cmd[i]
                yield snap(keyDelay(i) + (i === step.cmd.length - 1 ? 350 : 0))
            }
            line.cursor = false
        } else if ('ask' in step) {
            const answer: Seg = ['text', '']
            push([['yellow', `${step.ask} `], answer], true)
            yield snap(650)
            answer[1] = step.answer
            yield snap(250)
        } else if ('spin' in step) {
            const frame: Seg = ['yellow', FRAMES[0]]
            push([frame, ['text', ` ${step.spin}`]])
            for (let i = 0; i * 100 < step.ms; i++) {
                frame[1] = FRAMES[i % FRAMES.length]
                yield snap(100)
            }
            lines.pop()
        } else {
            push(step.out)
            yield snap(step.ms ?? (step.out.length ? 90 : 40))
        }
    }
    dir = END_DIR
    push(prompt(END_DIR), true)
    yield snap(HOLD)
}

/** The last frame of the story, with the cursor on a new prompt. */
export function finalLines(): Line[] {
    let last: Frame | undefined
    for (const frame of frames()) last = frame
    return last!.lines
}
