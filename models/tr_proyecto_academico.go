package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type TrProyectoAcademico struct {
	ProyectoAcademicoInstitucion *ProyectoAcademicoInstitucion
	Registro                     *[]RegistroCalificadoAcreditacion
	Enfasis                      *[]InstitucionEnfasis
	Titulaciones				 *[]Titulacion
}

type TrProyectoAcademicoPutInfoBasica struct {
	ProyectoAcademicoInstitucion *ProyectoAcademicoInstitucion
	Titulaciones				 *[]Titulacion
}

// GetProyectoAcademicasById Transacción para consultar todas las proyectos academicos con toda la información de las mismas
func GetProyectoAcademicasById(id int) (v []interface{}, err error) {
	fmt.Println("entro al modelo")
	o := orm.NewOrm()
	var proyectos []*ProyectoAcademicoInstitucion

	if _, err := o.QueryTable(new(ProyectoAcademicoInstitucion)).RelatedSel().Filter("Id", id).Filter("Activo", true).All(&proyectos); err == nil {

		for _, proyecto := range proyectos {

			proyectoAcademico := proyecto
			fmt.Println("err ojoooo /n", proyecto)

			var registroProyectos []RegistroCalificadoAcreditacion
			if _, err := o.QueryTable(new(RegistroCalificadoAcreditacion)).RelatedSel().Filter("Id", id).All(&registroProyectos); err != nil {
				//fmt.Println("registro/n", registroProyectos)
				return nil, err
			}

			var enfasiproyectos []InstitucionEnfasis
			if _, err := o.QueryTable(new(InstitucionEnfasis)).RelatedSel().Filter("Id", id).All(&enfasiproyectos); err != nil {
				return nil, err
			}

			v = append(v, map[string]interface{}{
				"Id":           proyectoAcademico.Id,
				"Codigo":       proyectoAcademico.Codigo,
				"Nombre":       proyectoAcademico.Nombre,
				"Codigo SNIES": proyectoAcademico.CodigoSnies,
				"Duracion":     proyectoAcademico.Duracion,
				"Año Acto":     proyectoAcademico.AnoActoAdministrativoId,
				"Registro":     &registroProyectos,
				"Enfasis":      &enfasiproyectos,
			})
		}

		return v, nil
	}
	return nil, err
}

// AddTransaccionProyectoAcademica Transacción para registrar toda la información de un proyecto academico
func AddTransaccionProyectoAcademica(m *TrProyectoAcademico) (err error) {
	o := orm.NewOrm()
	err = o.Begin()

	if idProyecto, errTr := o.Insert(m.ProyectoAcademicoInstitucion); errTr == nil {
		fmt.Println(idProyecto)

		for _, v := range *m.Registro {
			v.ProyectoAcademicoInstitucionId.Id = int(idProyecto)
			if _, errTr = o.Insert(&v); errTr != nil {
				err = errTr
				fmt.Println(err)
				_ = o.Rollback()
				return
			}
		}

		for _, v := range *m.Enfasis {
			v.ProyectoAcademicoInstitucionId.Id = int(idProyecto)
			if _, errTr = o.Insert(&v); errTr != nil {
				err = errTr
				fmt.Println(err)
				_ = o.Rollback()
				return
			}
		}

		for _, v := range *m.Titulaciones {
			v.ProyectoAcademicoInstitucionId.Id = int(idProyecto)
			if _, errTr = o.Insert(&v); errTr != nil {
				err = errTr
				fmt.Println(err)
				_ = o.Rollback()
				return
			}
		}

		_ = o.Commit()
	} else {
		err = errTr
		fmt.Println(err)
		_ = o.Rollback()
	}

	return
}

// UpdateTransaccionProyectoAcademica Transacción para actualizar toda la información básicoa de un proyecto academico
func UpdateTransaccionProyectoAcademica(m *TrProyectoAcademicoPutInfoBasica) (err error) {
	o := orm.NewOrm()
	err = o.Begin()

	//Actualizar el proyecto academico
	if idProyecto, errTr := o.Update(m.ProyectoAcademicoInstitucion); errTr == nil {
		fmt.Println(idProyecto)

		for _, v := range *m.Titulaciones {
			var titulacion Titulacion
			fmt.Println("TipoTitulacionId__Id",v.TipoTitulacionId.Id,"ProyectoAcademicoInstitucionId__Id",v.ProyectoAcademicoInstitucionId.Id)
			if errTr := o.QueryTable(new(Titulacion)).RelatedSel().Filter("TipoTitulacionId__Id",v.TipoTitulacionId.Id).Filter("ProyectoAcademicoInstitucionId__Id",v.ProyectoAcademicoInstitucionId.Id).One(&titulacion); errTr == nil{
				// Si existe hace update
				v.Id = titulacion.Id
				if _, errTr = o.Update(&v,"Nombre","Descripcion","Activo","NumeroOrden"); errTr != nil {
					err = errTr
					fmt.Println(err)
					_ = o.Rollback()
					return
				} else {
					fmt.Println("update para titulacion",titulacion.Id)
				}
			} else {
				if errTr == orm.ErrNoRows {
					if idTitulacion, errTr := o.Insert(&v); errTr != nil {
						err = errTr
						fmt.Println(err)
						_ = o.Rollback()
						return
					} else {
						fmt.Println("insert para inscripcion",idTitulacion)
					}
				} else {
					err = errTr
					fmt.Println(err)
					_ = o.Rollback()
					return
				}
			}
		}
		/* Transsacciones a parte
		for _, v := range *m.Registro {
			// v.ProyectoAcademicoInstitucionId.Id = int(idProyecto)
			if _, errTr = o.Update(&v); errTr != nil {
				err = errTr
				fmt.Println(err)
				_ = o.Rollback()
				return
			}
		}
		
		for _, v := range *m.Enfasis {
			// v.ProyectoAcademicoInstitucionId.Id = int(idProyecto)
			if _, errTr = o.Update(&v); errTr != nil {
				err = errTr
				fmt.Println(err)
				_ = o.Rollback()
				return
			}
		}
		*/
		_ = o.Commit()
	} else {
		err = errTr
		fmt.Println(err)
		_ = o.Rollback()
	}

	return
}