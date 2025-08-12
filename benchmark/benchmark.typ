#import "lib.typ": tpl

// ================================================
// LAYOUT & STYLE
// ================================================

#show: tpl.with(
)

// ================================================
// METHODS
// ================================================

#let pangram = "Nyx’ Bö drückt Vamps Quiz-Floß jäh weg."

// Measurement tape scale with cm and mm markings
#let measure_tape() = {
  let total_width = 79.5mm // slightly less than page width for margin
  let cm_count = 8 // number of cm marks to show

  // Main ruler line
  line(length: total_width, stroke: 0.5pt)

  // Generate cm and mm marks
  for cm in range(cm_count) {
    let x_pos = cm * 10mm

    // CM mark (longer line)
    place(top + left, dx: x_pos, dy: 0mm, line(angle: 90deg, length: 3mm, stroke: 0.5pt))

    // CM number
    place(top + left, dx: x_pos + 1mm, dy: 3.5mm, text(size: 6pt, str(cm)))

    // MM marks (shorter lines)
    for mm in range(1, 10) {
      let mm_x = x_pos + mm * 1mm
      if mm_x < total_width {
        let mark_height = if mm == 5 { 2mm } else { 1mm }
        place(top + left, dx: mm_x, dy: 0mm, line(angle: 90deg, length: mark_height, stroke: 0.3pt))
      }
    }
  }
}

// Function to create a hatched rectangle
#let hatched_rect(
  size: 4.0mm,
  hatch_type: "vertical", // "vertical" or "diagonal"
  hatch_count: 10,
  hatch_thickness: 0.5pt,
  hatch_color: black,
  label: none, // optional label text above the rectangle
) = box(width: size, height: size)[
  // Optional label above the rectangle
  #if label != none {
    align(center, text(size: 5pt, label))
  }

  // Create the base rectangle (transparent background)
  #rect(
    width: size,
    height: size,
    fill: none,
    stroke: none,
  )

  // Add hatching lines
  #for i in range(hatch_count + 1) {
    let spacing = size / hatch_count
    let x_offset = i * spacing

    if hatch_type == "vertical" {
      // Vertical hatching
      place(
        top + left,
        dx: x_offset,
        dy: if label != none { 8pt } else { 0mm },
        line(
          angle: 90deg,
          length: size,
          stroke: hatch_thickness + hatch_color,
        ),
      )
    } else if hatch_type == "horizontal" {}
  }
]

#let hatched_rect_diagonal(
  size: 4.0mm,
  hatch_count: 10,
  hatch_thickness: 0.5pt,
  hatch_color: black,
  label: none, // optional label text above the rectangle
) = box(width: size, height: size)[
  // Calculate diagonal lines for hatching
  #let spacing = size / (hatch_count + 1)

  #for i in range(hatch_count) {
    let offset = (i + 1) * spacing

    // First set of diagonal lines (from top-left to bottom-right direction)
    // Line starts from left edge and goes diagonally
    if offset <= size {
      place(
        top + left,
        line(
          start: (0mm, offset),
          end: (offset, 0mm),
          stroke: hatch_thickness + hatch_color,
        ),
      )
    }

    // Second set of diagonal lines (continuing the pattern)
    // Line starts from bottom edge and goes diagonally
    if offset <= size {
      place(
        top + left,
        line(
          start: (offset, size),
          end: (size, offset),
          stroke: hatch_thickness + hatch_color,
        ),
      )
    }
  }

  // Add the main diagonal line from bottom-left to top-right
  #place(
    top + left,
    line(
      start: (0mm, size),
      end: (size, 0mm),
      stroke: hatch_thickness + hatch_color,
    ),
  )
]

#let hatch_table_row(thickness, counts) = (
  rotate(-90deg, [#thickness.pt()]),
  ..counts.map(it => hatched_rect(hatch_count: it, hatch_thickness: thickness)),
  rotate(-90deg, [#thickness.pt()]),
  ..counts.map(it => hatched_rect_diagonal(hatch_count: it, hatch_thickness: thickness)),
)


// ================================================
// CONTENT
// ================================================

// Display the measurement tape
#measure_tape()
#v(5mm)


#columns[
  = Lines
  #{
    let thicknesses = (0.25pt, 0.5pt, 0.75pt, 1pt, 1.5pt, 2pt)
    set text(size: 4.5pt)
    grid(
      columns: (auto, 1fr),
      inset: (
        top: 0.5mm,
        left: 1mm,
        right: 1mm,
        bottom: 0.5mm,
      ),
      ..for thickness in thicknesses {
        (
          [#thickness.pt()],
          grid.cell(
            stroke: (bottom: thickness),
            [],
          ),
        )
      }
    )
  }
  = Text Size
  #for size in (4pt, 5pt, 6pt, 7pt, 8pt, 9pt, 10pt) {
    set text(size: size)
    [ABCabc123 (#size.pt()#sym.space.hair\pt)]
    linebreak()
  }

  = Gray Scale

  #{
    let steps = 16
    let step_size = 255 / (steps - 1)
    set text(size: 6.0pt)
    grid(
      columns: range(4).map(_ => 1fr),
      rows: 4mm,
      gutter: 1mm,
      align: center + horizon,
      ..range(steps).map(it => grid.cell(
        fill: luma(calc.floor(it * step_size)),
        [
          #let density = calc.floor(it * step_size)
          #set text(fill: white) if density < 128
          #density
        ],
      ))
    )
  }

  #rect(width: 100%, height: 4mm, fill: gradient.linear(white, black))

  = Text (6.5pt)
  #text(size: 6.5pt, pangram)

  #colbreak()

  = Hatching

  #{
    let counts = (6, 8, 10, 12, 14, 16)
    set text(size: 4.5pt)
    table(
      columns: (3mm, ..counts.map(_ => auto)),
      inset: (
        top: 0.5mm,
        left: 1mm,
        right: 1mm,
        bottom: 0.5mm,
      ),
      stroke: none,
      align: center + horizon,
      table.header([], ..counts.map(it => [#it])),

      ..hatch_table_row(0.05pt, counts),
      ..hatch_table_row(0.1pt, counts),
      ..hatch_table_row(0.3pt, counts),
      ..hatch_table_row(0.5pt, counts),
    )
  }

  = Fonts
  #{
    let fonts = ("Arial", "Courier New", "Times New Roman", "Verdana", "Palatino", "Comic Sans MS", "Impact")
    set text(size: 6.0pt)
    set par(spacing: 0.6em)
    for ft in fonts [
      #ft.slice(0, 2):
      #set text(font: ft)
      #pangram

    ]
  }
]

#image("cave.jpg")