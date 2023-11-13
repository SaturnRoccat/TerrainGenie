
export class Chunk {
    constructor(xpos, zpos, xsize, ysize, blocks) {
        this.x = xpos;
        this.z = zpos;
        this.xsize = xsize;
        this.ysize = ysize;
        this.zsize = xsize;
        this.blocks = blocks;
    }

    getBlock(x, y, z) {
        return this.blocks[x + y * this.xsize + z * this.xsize * this.ysize];
    }
}