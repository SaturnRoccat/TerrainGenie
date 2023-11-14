export class Chunk {
    constructor(xPos, zPos, xSize, ySize, blocks) {
        this.x = xPos;
        this.z = zPos;
        this.xSize = xSize;
        this.ySize = ySize;
        this.zSize = xSize;
        this.blocks = blocks;
    };

    getBlock(x, y, z) {
        return (
            this.blocks[
                x + ( y * this.xSize ) + ( z * this.xSize * this.ySize )
            ]
        );
    };
};