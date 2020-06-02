package xapi

import (
	"github.com/arturoguerra/xcpng-csi/internal/structs"
	xenapi "github.com/terra-farm/go-xen-api-client"
)

func (c *xClient) IsAttached(volID, nodeID string, zone *structs.Zone) (bool, error) {
	api, session, err := c.Connect(zone)
	if err != nil {
		return false, err
	}
	defer c.Close(api, session)

	vm, err := c.GetVM(api, session, nodeID)
	if err != nil {
		return false, err
	}

	log.Info("VDI.GetAllRecords")
	vdis, err := api.VDI.GetAllRecords(session)
	if err != nil {
		return false, err
	}

	var vdiUUID xenapi.VDIRef
	for ref, vdi := range vdis {
		if vdi.NameLabel == volID && !vdi.IsASnapshot {
			vdiUUID = ref
		}
	}

	log.Info("VBD.GetAllRecords")
	vbds, err := api.VBD.GetAllRecords(session)
	if err != nil {
		return false, err
	}

	for _, vbd := range vbds {
		if vbd.VM == vm && vbd.CurrentlyAttached && vbd.VDI == vdiUUID {
			return true, nil
		}
	}

	return false, nil
}
